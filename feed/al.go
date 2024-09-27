package feed

import (
	"context"
	"errors"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/gofiber/fiber/v2"
	"github.com/watchedsky-social/core/internal/config"
	"github.com/watchedsky-social/core/internal/utils"
)

type FeedGeneratorInput struct {
	Limit        uint
	Cursor       *string
	RequestorDID *string
}

type FeedAlgorithm interface {
	ShortName() string
	UserSpecific() bool
	GenerateFeed(context.Context, FeedGeneratorInput) (*bsky.FeedGetFeedSkeleton_Output, error)
}

var registeredFeeds []FeedAlgorithm = []FeedAlgorithm{}

func handleFeedAlgorithm(c *fiber.Ctx) error {
	feedCfg := config.Config.Feed

	feed := c.Query("feed")
	limit := c.QueryInt("limit", 30)
	cursor := c.Query("cursor")

	atURI, err := syntax.ParseATURI(feed)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}

	authority := atURI.Authority().String()
	if authority != feedCfg.PublisherDID ||
		atURI.Collection().String() != "app.bsky.feed.generator" {
		return c.SendStatus(http.StatusNotFound)
	}

	requestedAlgo := atURI.RecordKey().String()
	idx := slices.IndexFunc(registeredFeeds, func(a FeedAlgorithm) bool {
		return a.ShortName() == requestedAlgo
	})

	if idx == -1 {
		return c.SendStatus(http.StatusNotFound)
	}

	algo := registeredFeeds[idx]
	requestorDID := ""
	if algo.UserSpecific() {
		authHeader := http.Header(c.GetReqHeaders()).Get("Authorization")
		if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			log.Printf("auth header %q", authHeader)
			return c.SendStatus(http.StatusUnauthorized)
		}

		jwt := authHeader[7:]
		claims, err := verifyAuth(c.UserContext(), jwt, feedCfg.ServiceDID, "app.bsky.feed.generator")
		if err != nil {
			log.Printf("err from verify auth: %v", err)
			return c.SendStatus(http.StatusUnauthorized)
		}

		iss, err := claims.GetIssuer()
		if err != nil {
			return c.SendStatus(http.StatusForbidden)
		}

		requestorDID = iss
	}

	input := FeedGeneratorInput{
		Limit:        uint(limit),
		Cursor:       utils.NilRefIfZero(cursor),
		RequestorDID: utils.NilRefIfZero(requestorDID),
	}

	output, err := algo.GenerateFeed(c.UserContext(), input)
	if err != nil {
		log.Println(err)
		var f FeedHTTPError

		if errors.As(err, &f) {
			return c.SendStatus(f.StatusCode())
		}

		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.JSON(output)
}
