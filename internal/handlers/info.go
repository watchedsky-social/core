package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/watchedsky-social/core/internal/config"
)

var buildInfo = map[string]string{
	"version":  config.Version,
	"build_id": config.BuildID,
}

func Info(ctx *fiber.Ctx) error {
	return ctx.JSON(buildInfo)
}
