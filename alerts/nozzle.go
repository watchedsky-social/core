package alerts

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bluesky-social/indigo/api"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/bluesky-social/indigo/did"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/watchedsky-social/core/internal/database/models"
	"github.com/watchedsky-social/core/internal/database/query"
	"gorm.io/gorm/clause"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/events"
	"github.com/bluesky-social/indigo/events/schedulers/sequential"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/repomgr"
	"github.com/gorilla/websocket"
	"github.com/ipfs/go-cid"
)

const (
	firehoseURL = "wss://bsky.network/xrpc/com.atproto.sync.subscribeRepos"
	marker      = "🌩️👀 "
)

type Request struct {
	CID string
	DID string
	WID string
}

func SubscribeToFirehose(ctx context.Context) error {
	targetDID := os.Getenv("FIREHOSE_TARGET_DID")
	if targetDID == "" {
		return errors.New("env variable FIREHOSE_TARGET_DID must be set")
	}

	hr, err := api.NewProdHandleResolver(1, "https://plc.directory", false)
	if err != nil {
		return fmt.Errorf("could not create handle resolver: %w", err)
	}

	targetHandle, _, err := api.ResolveDidToHandle(ctx, did.NewMultiResolver(), hr, targetDID)
	if err != nil {
		return fmt.Errorf("could not resolve did %q: %w", targetDID, err)
	}

	con, _, err := websocket.DefaultDialer.Dial(firehoseURL, http.Header{})
	if err != nil {
		return fmt.Errorf("error connecting to firehose: %w", err)
	}
	defer con.Close()

	go func() {
		<-ctx.Done()
		con.Close()
	}()

	xrpcClient, err := getAuthenticatedXRPCClient(ctx)
	if err != nil {
		return fmt.Errorf("cannot connect to bluesky: %w", err)
	}

	firehoseNozzle := &events.RepoStreamCallbacks{
		RepoCommit: func(evt *atproto.SyncSubscribeRepos_Commit) error {
			if evt.TooBig {
				log.Printf("skipping too big events for now: %d", evt.Seq)
				return nil
			}

			r, err := repo.ReadRepoFromCar(ctx, bytes.NewReader(evt.Blocks))
			if err != nil {
				log.Printf("reading repo from car (seq: %d, len: %d): %v", evt.Seq, len(evt.Blocks), err)
				return nil
			}

			for _, op := range evt.Ops {
				ek := repomgr.EventKind(op.Action)
				switch ek {
				case repomgr.EvtKindCreateRecord, repomgr.EvtKindUpdateRecord:
					rc, rec, err := r.GetRecord(ctx, op.Path)
					if err != nil {
						e := fmt.Errorf("getting record %s (%s) within seq %d for %s: %w", op.Path, *op.Cid, evt.Seq, evt.Repo, err)
						log.Print(e)
						continue
					}

					if lexutil.LexLink(rc) != *op.Cid {
						// TODO: do we even error here?
						log.Printf("mismatch in record and op cid: %s != %s", rc, *op.Cid)
						continue
					}

					if err := handleCreateRecord(ctx, evt.Repo, &rc, rec, targetDID, targetHandle, xrpcClient); err != nil {
						log.Printf("event consumer callback (%s): %s", ek, err)
						continue
					}
				}
			}

			return nil
		},
	}

	return events.HandleRepoStream(ctx, con, sequential.NewScheduler("watchedsky-nozzle", firehoseNozzle.EventHandler))
}

func handleCreateRecord(ctx context.Context, did string, rcid *cid.Cid, rec any, targetDID string, targetHandle string, xc *xrpc.Client) error {
	orig, isPost := rec.(*bsky.FeedPost)
	if !isPost {
		return nil
	}

	markerIDX := strings.Index(orig.Text, marker)
	// either the marker wasn't found or it wasn't followed by a 12 character watch ID
	if markerIDX == -1 || markerIDX > len(orig.Text)-(len(marker)+12) {
		return nil
	}

	watchID := orig.Text[markerIDX+len(marker) : markerIDX+len(marker)+12]

	log.Printf("found potential watch ID %q in post %s by did %s", watchID, rcid, did)

	log.Println("Look to see if they @'ed us...")
	if !strings.HasPrefix(targetHandle, "@") {
		targetHandle = "@" + targetHandle
	}

	saveAndRespond := false
	if strings.Contains(orig.Text, targetHandle) {
		// let's see if there's a facet there to double check
		for _, facet := range orig.Facets {
			for _, feat := range facet.Features {
				if feat.RichtextFacet_Mention != nil && feat.RichtextFacet_Mention.Did == targetDID {
					saveAndRespond = true
					break
				}
			}
		}
	}

	var root *atproto.RepoStrongRef = nil

	// maybe they just replied to us! we should check that too. Don't go too crazy,
	// only check direct replies.
	if !saveAndRespond && orig.Reply != nil && orig.Reply.Parent != nil {
		root = orig.Reply.Root
		uri, err := syntax.ParseATURI(orig.Reply.Parent.Uri)
		if err != nil {
			return fmt.Errorf("error parsing at:// URI %s: %v", orig.Reply.Parent.Uri, err)
		}

		did, err := uri.Authority().AsDID()
		if err != nil {
			return fmt.Errorf("at URI has no authority: %v", err)
		}

		saveAndRespond = did.String() == targetDID
	}

	if saveAndRespond {
		// first, save
		dao := query.Didwid.WithContext(ctx)
		didwid := &models.Didwid{
			Did: did,
			Wid: watchID,
		}

		if err := dao.Clauses(clause.OnConflict{UpdateAll: true}).Create(didwid); err != nil {
			return fmt.Errorf("couldn't update didwids: %w", err)
		}

		// now, respond. if the root of the thread is nil, then it's the original post, which is
		// also the parent
		parentCID := rcid.String()
		parentURI := fmt.Sprintf("at://%s/%s/%s", did, orig.LexiconTypeID, rcid.String())
		parent := &atproto.RepoStrongRef{
			Cid: parentCID,
			Uri: parentURI,
		}

		if root == nil {
			root = parent
		}

		input := &bsky.FeedPost{
			Text:      "Hey there 👋🏻\n\nI saved your watch ID, so you should be able to see your updated feed now. If you want, you can delete your skeet and I won't mind. Just @ me or reply again to change your watch ID!",
			Langs:     []string{"en-US"},
			CreatedAt: time.Now().Format(time.RFC3339),
			Reply: &bsky.FeedPost_ReplyRef{
				Root:   root,
				Parent: parent,
			},
		}

		_, err := atproto.RepoCreateRecord(ctx, xc, &atproto.RepoCreateRecord_Input{
			Collection: "app.bsky.feed.post",
			Repo:       xc.Auth.Did,
			Record: &lexutil.LexiconTypeDecoder{
				Val: input,
			},
		})

		return err
	}

	return nil
}
