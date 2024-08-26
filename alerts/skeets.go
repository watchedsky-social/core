package alerts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/util"

	"github.com/bluesky-social/indigo/xrpc"
	"github.com/watchedsky-social/core/internal/config"
	"github.com/watchedsky-social/core/internal/database/models"
	"github.com/watchedsky-social/core/internal/database/query"
	"github.com/watchedsky-social/core/internal/utils"
	"gorm.io/gorm"
)

type alertsState struct {
	AllAlerts   map[string]*models.Alert           `json:"allAlerts"`
	NewAlertIDs []string                           `json:"newAlertIDs"`
	ThreadMap   map[string]*bsky.FeedPost_ReplyRef `json:"threadMap"`
}

const (
	stateFile    = "state.json"
	skeetRate    = time.Hour / 1666 // current rate limit
	baseURL      = "https://watchedsky.social/app/alerts"
	skeetSuffix  = " See more: watchedsky.social/app/alerts/..." // leading space is intentional
	skeetFormat  = "NEW WEATHER ADVISORY: %s\n\n%s"
	maxSkeetLen  = 300
	shortLinkLen = len("watchedsky.social/app/alerts/...")
)

func hydrateAlertState(ctx context.Context) (err error) {
	output := alertsState{
		AllAlerts:   make(map[string]*models.Alert),
		NewAlertIDs: make([]string, 0, 10),
		ThreadMap:   make(map[string]*bsky.FeedPost_ReplyRef),
	}
	defer func() {
		data, e := json.Marshal(output)
		if e != nil {
			err = e
			return
		}

		bskyCfg := config.Config.Alerts
		e2 := os.WriteFile(filepath.Join(bskyCfg.HydrationDir, stateFile), data, 0o666)
		if err != nil || e2 != nil {
			errs := utils.Filter([]error{err, e2}, func(e error) bool { return e != nil })
			err = errors.Join(errs...)
			return
		}
	}()

	alert := query.Alert
	dao := alert.WithContext(ctx)
	var results []*models.Alert
	results, err = dao.Select(
		alert.ID, alert.Headline, alert.Description, alert.Sent,
		alert.ReferenceIds, alert.SkeetInfo, alert.MessageType,
	).Order(alert.Sent.Asc()).Find()

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}

		return
	}

	safeResults := utils.Filter(results, func(a *models.Alert) bool { return a != nil })
	for _, alert := range safeResults {
		output.AllAlerts[alert.ID] = alert
		if alert.SkeetInfo == nil {
			output.NewAlertIDs = append(output.NewAlertIDs, alert.ID)
		}
		if alert.ReferenceIds != nil && len(*alert.ReferenceIds) > 0 {
			output.ThreadMap[alert.ID] = &bsky.FeedPost_ReplyRef{}
		}
	}

	for alertID := range output.ThreadMap {
		alert := output.AllAlerts[alertID]
		refIDs := *alert.ReferenceIds

		parentID := refIDs[0]
		rootID := refIDs[0]
		if len(refIDs) > 1 {
			rootID = refIDs[len(refIDs)-1]
		}

		parentAlert := output.AllAlerts[parentID]
		rootAlert := output.AllAlerts[rootID]
		if rootAlert == nil {
			rootAlert = parentAlert
		}

		if parentAlert != nil {
			output.ThreadMap[alertID] = &bsky.FeedPost_ReplyRef{
				Parent: &atproto.RepoStrongRef{
					Uri: parentAlert.SkeetInfo.Uri,
					Cid: parentAlert.SkeetInfo.Cid,
				},
				Root: &atproto.RepoStrongRef{
					Uri: rootAlert.SkeetInfo.Uri,
					Cid: rootAlert.SkeetInfo.Cid,
				},
			}
		}
	}

	return
}

func InitialHydration(ctx context.Context) error {
	log.Println("Checking the radar map...")

	return hydrateAlertState(ctx)
}

func SkeetNewAlerts(ctx context.Context) error {
	log.Println("Posting about new alerts...")
	cfg := config.Config.Alerts

	var state alertsState
	stateFile, err := os.Open(filepath.Join(cfg.HydrationDir, stateFile))
	if err != nil {
		return err
	}
	defer stateFile.Close()

	if err = json.NewDecoder(stateFile).Decode(&state); err != nil {
		return err
	}

	xrpcClient := &xrpc.Client{
		Client:    util.RobustHTTPClient(),
		Host:      "https://bsky.social",
		UserAgent: utils.Ref("watchedsky.social"),
	}

	session, err := atproto.ServerCreateSession(ctx, xrpcClient, &atproto.ServerCreateSession_Input{
		Identifier: cfg.Bluesky.ID,
		Password:   cfg.Bluesky.AppPassword,
	})

	if err != nil {
		return err
	}

	xrpcClient.Auth = &xrpc.AuthInfo{
		AccessJwt:  session.AccessJwt,
		RefreshJwt: session.RefreshJwt,
		Handle:     session.Handle,
		Did:        session.Did,
	}

	alertErrors := []error{}
	for _, alertID := range state.NewAlertIDs {
		if err = postSkeet(ctx, alertID, xrpcClient, &state); err != nil {
			alertErrors = append(alertErrors, fmt.Errorf("alert %s could not be posted to Bluesky: %w", alertID, err))
		}

		time.Sleep(skeetRate)
	}

	if len(alertErrors) > 0 {
		log.Printf("WARNING: the following errors occurred while posting alerts to Bluesky: %v", errors.Join(alertErrors...))
	}

	return refreshDB(ctx, state)
}

func postSkeet(ctx context.Context, alertID string, xrpcClient *xrpc.Client, state *alertsState) error {
	alert, ok := state.AllAlerts[alertID]
	if !ok {
		return fmt.Errorf("alert %s not found", alertID)
	}

	alertURL := fmt.Sprintf("%s/%s", baseURL, alertID)

	cleanHeadline := alert.Headline
	cutoff := strings.Index(strings.ToLower(alert.Headline), "issued")
	if cutoff > 0 {
		cleanHeadline = cleanHeadline[0:cutoff]
	}
	cleanDescription := strings.Join(strings.Fields(alert.Description), " ")
	fullSkeetBody := strings.TrimSpace(fmt.Sprintf(skeetFormat, cleanHeadline, cleanDescription))

	fullSkeetLen := len(fullSkeetBody)

	skeetLen := min(maxSkeetLen-len(skeetSuffix), fullSkeetLen)
	fullSkeetBody = fullSkeetBody[0:skeetLen]

	isTruncated := fullSkeetLen > skeetLen
	if isTruncated {
		lastSpace := strings.LastIndexFunc(fullSkeetBody, func(r rune) bool { return r == ' ' })
		fullSkeetBody = fullSkeetBody[0:lastSpace] + "..."
	}

	fullSkeetBody = fullSkeetBody + skeetSuffix
	input := &bsky.FeedPost{
		CreatedAt: alert.Sent.Format(time.RFC3339),
		Langs:     []string{"en_US"},
		Text:      fullSkeetBody,
		Embed: &bsky.FeedPost_Embed{
			EmbedExternal: &bsky.EmbedExternal{
				External: &bsky.EmbedExternal_External{
					Title:       "WatchedSky",
					Description: "National Weather Service Alerts, now on Bluesky",
					Uri:         alertURL,
				},
			},
		},
		Facets: []*bsky.RichtextFacet{
			{
				Index: &bsky.RichtextFacet_ByteSlice{
					ByteStart: int64(len(fullSkeetBody) - shortLinkLen),
					ByteEnd:   int64(len(fullSkeetBody) - 1),
				},
				Features: []*bsky.RichtextFacet_Features_Elem{
					{
						RichtextFacet_Link: &bsky.RichtextFacet_Link{
							Uri: alertURL,
						},
					},
				},
			},
		},
		Reply: state.ThreadMap[alertID],
	}

	response, err := atproto.RepoCreateRecord(ctx, xrpcClient, &atproto.RepoCreateRecord_Input{
		Record: &lexutil.LexiconTypeDecoder{
			Val: input,
		},
		Collection: "app.bsky.feed.post",
		Repo:       xrpcClient.Auth.Handle,
	})

	if err != nil {
		return err
	}

	alert.SkeetInfo = &models.SkeetInfo{
		Cid: response.Cid,
		Uri: response.Uri,
	}

	state.AllAlerts[alertID] = alert
	return nil
}

func refreshDB(ctx context.Context, state alertsState) error {
	return query.Q.Transaction(func(tx *query.Query) error {
		dao := tx.Alert
		for _, id := range state.NewAlertIDs {
			alert := state.AllAlerts[id]
			if alert == nil {
				return fmt.Errorf("no alert with id %s", id)
			}

			if alert.SkeetInfo != nil {
				if _, err := dao.WithContext(ctx).
					Where(dao.ID.Eq(id)).
					UpdateSimple(dao.SkeetInfo.Value(alert.SkeetInfo)); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
