package api

import (
	"context"
	"fmt"
	"guigou/bot-is-gud/api/rpc"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"go.uber.org/zap"
)

type API struct {
	logger      *zap.SugaredLogger
	startedAt   time.Time
	lastTypedAt *time.Time
	protoClient rpc.Presence
	isHealthy   isHealthy
	ctx         context.Context
}

type isHealthy func() bool

func New(logger *zap.SugaredLogger, port string, lastTypedAt *time.Time, ctx context.Context) *API {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	client := rpc.NewPresenceProtobufClient("http://localhost:8080", &http.Client{})
	api := &API{logger: logger, startedAt: time.Now().UTC(), lastTypedAt: lastTypedAt, protoClient: client, isHealthy: func() bool { return true }, ctx: ctx}
	app.Use(limiter.New())
	app.Get("/", api.healthHandler)
	app.Get("/whoseOn", api.whoseOnHandler)
	go app.Listen(":" + port)
	return api
}

func (api *API) RegisterHealthCheck(fn isHealthy) {
	api.isHealthy = fn
}

func (api *API) healthHandler(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	statusCode := 200
	statusMessage := "HEALTHY"
	if !api.isHealthy() {
		statusCode = 503
		statusMessage = "UNHEALTHY"
	}
	return c.Status(statusCode).JSON(&fiber.Map{
		"uptime":      fmt.Sprintf("%s", time.Now().UTC().Sub(api.startedAt).Round(time.Second)),
		"lastTypedAt": api.lastTypedAt.Format(time.RFC3339),
		"status":      statusMessage,
	})
}

func (api *API) whoseOnHandler(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	resp, err := api.protoClient.WhoseOn(api.ctx, &rpc.WhoseOnReq{VoiceChannel: ""})
	if err != nil {
		api.logger.Warn("unable to fetch whose on: %s", err)
		return c.Status(500).JSON(&fiber.Map{
			"error": "Internal Server Error",
		})
	}

	return c.Status(200).JSON(&fiber.Map{
		"count": len(resp.Users),
	})
}
