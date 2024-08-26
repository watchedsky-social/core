package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/watchedsky-social/core/internal/database/models"
	"github.com/watchedsky-social/core/internal/database/query"
	"github.com/watchedsky-social/core/internal/utils"
)

type recentAlert struct {
	ID          string `json:"id"`
	Headline    string `json:"headline"`
	Description string `json:"description"`
	Sent        int64  `json:"sent"`
	AreaDesc    string `json:"area_desc"`
}

func RecentAlerts(ctx *fiber.Ctx) error {
	dao := query.Alert
	alerts, err := dao.WithContext(ctx.UserContext()).
		Select(dao.ID, dao.Headline, dao.Description, dao.Sent, dao.AreaDesc).
		Order(dao.Sent.Desc()).
		Limit(ctx.QueryInt("lim", 20)).
		Find()

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return err
	}

	return ctx.JSON(utils.Map(alerts, func(a *models.Alert) recentAlert {
		return recentAlert{
			ID:          a.ID,
			Headline:    a.Headline,
			Description: a.Description,
			Sent:        a.Sent.UTC().UnixMilli(),
			AreaDesc:    a.AreaDesc,
		}
	}))
}

func AlertByID(ctx *fiber.Ctx) error {
	dao := query.Alert
	alert, err := dao.WithContext(ctx.UserContext()).Where(dao.ID.Eq(ctx.Params("alertid"))).Find()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return err
	}

	if len(alert) == 0 {
		ctx.Status(http.StatusNotFound)
		return nil
	}

	return ctx.JSON(alert[0])
}
