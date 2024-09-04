package feed

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/util"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/watchedsky-social/core/internal/database/models"
	"github.com/watchedsky-social/core/internal/database/query"
	"github.com/watchedsky-social/core/internal/utils"
)

type localAlerts struct {
	xc *xrpc.Client
}

func newLocalAlerts() *localAlerts {
	return &localAlerts{
		xc: &xrpc.Client{
			Client: util.RobustHTTPClient(),
			Host:   "https://public.api.bsky.app",
		},
	}
}

func transformAlerts(alerts []*models.Alert) []*bsky.FeedDefs_SkeletonFeedPost {
	return utils.Map(alerts, func(alert *models.Alert) *bsky.FeedDefs_SkeletonFeedPost {
		return &bsky.FeedDefs_SkeletonFeedPost{
			Post: alert.SkeetInfo.Uri,
		}
	})
}

const tag = "🌩️👀 "

func (l *localAlerts) getWID(ctx context.Context, did string) (string, error) {
	didwid, err := query.Didwid.WithContext(ctx).Select(query.Didwid.Wid).Where(query.Didwid.Did.Eq(did)).First()
	if err == nil {
		return didwid.Wid, nil
	}

	out, err := bsky.ActorGetProfile(ctx, l.xc, did)
	if err != nil {
		return "", err
	}

	if out.Description != nil {
		desc := *out.Description
		idx := strings.Index(desc, tag)
		if idx != -1 && idx < len(desc)-len(tag)-13 {
			wid := desc[idx+len(tag)+1 : idx+len(tag)+13]
			if err = query.Didwid.WithContext(ctx).Create(&models.Didwid{
				Did: did,
				Wid: wid,
			}); err != nil {
				// don't fail here just because we couldn't save, the wid is in a stable place
				log.Println(err)
				return wid, nil
			}
		}
	}

	return "", errors.New("could not find wid")
}

// GenerateFeed implements FeedAlgorithm.
func (l *localAlerts) GenerateFeed(ctx context.Context, params FeedGeneratorInput) (*bsky.FeedGetFeedSkeleton_Output, error) {
	if params.RequestorDID == nil {
		return nil, NewFeedHTTPError(http.StatusUnauthorized, "Missing Authorization header")
	}

	watchID, err := l.getWID(ctx, *params.RequestorDID)
	if err != nil {
		log.Println(err)
		return &bsky.FeedGetFeedSkeleton_Output{
			Feed: []*bsky.FeedDefs_SkeletonFeedPost{
				{
					Post: noWatchIDFound,
				},
			},
		}, nil
	}

	dao := query.Alert.WithContext(ctx)
	var alerts []*models.Alert
	if params.Cursor == nil {
		alerts, err = dao.GetCustomAlertURIs(watchID, params.Limit)
	} else {
		var cursor uint64
		cursor, err = strconv.ParseUint(*params.Cursor, 10, 64)
		if err != nil {
			log.Println(err)
			return &bsky.FeedGetFeedSkeleton_Output{
				Feed: []*bsky.FeedDefs_SkeletonFeedPost{
					{
						Post: errorOccured,
					},
				},
			}, nil
		}
		alerts, err = dao.GetCustomAlertURIsWithCursor(watchID, params.Limit, uint(cursor))
	}

	if err != nil {
		log.Println(err)
		return &bsky.FeedGetFeedSkeleton_Output{
			Feed: []*bsky.FeedDefs_SkeletonFeedPost{
				{
					Post: errorOccured,
				},
			},
		}, nil
	}

	if len(alerts) == 0 {
		return &bsky.FeedGetFeedSkeleton_Output{
			Feed: []*bsky.FeedDefs_SkeletonFeedPost{
				{
					Post: noAlertsFound,
				},
			},
		}, nil
	}

	newCursor := ""
	if len(alerts) >= int(params.Limit) {
		oldestAlert := alerts[len(alerts)-1]
		newCursor = fmt.Sprint(oldestAlert.Sent.UnixMilli())
	}

	return &bsky.FeedGetFeedSkeleton_Output{
		Feed:   transformAlerts(alerts),
		Cursor: utils.NilRefIfZero(newCursor),
	}, nil
}

// ShortName implements FeedAlgorithm.
func (l *localAlerts) ShortName() string {
	return "watchedsky"
}

// UserSpecific implements FeedAlgorithm.
func (l *localAlerts) UserSpecific() bool {
	return true
}

func init() {
	registeredFeeds = append(registeredFeeds, newLocalAlerts())
}
