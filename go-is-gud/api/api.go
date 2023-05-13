package api

import (
	"context"
	"fmt"
	"guigou/bot-is-gud/api/rpc"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type API struct {
	startedAt   time.Time
	lastTypedAt *time.Time
	protoClient rpc.Presence
}

func New(lastTypedAt *time.Time, port string) {
	app := fiber.New()
	client := rpc.NewPresenceProtobufClient("http://localhost:8080", &http.Client{})
	api := &API{startedAt: time.Now().UTC(), lastTypedAt: lastTypedAt, protoClient: client}
	app.Use(limiter.New())
	app.Get("/", api.healthHandler)
	app.Get("/whoseOn", api.whoseOnHandler)
	app.Listen(":" + port)
}

func (api *API) healthHandler(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(200).JSON(&fiber.Map{
		"uptime":      fmt.Sprintf("%s", time.Now().UTC().Sub(api.startedAt).Round(time.Second)),
		"lastTypedAt": api.lastTypedAt.Format(time.RFC3339),
		"status":      "Ok",
	})
}

func (api *API) whoseOnHandler(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	resp, err := api.protoClient.WhoseOn(context.Background(), &rpc.WhoseOnReq{VoiceChannel: ""})
	if err != nil {
		fmt.Printf("unable to fetch whose on: %s\n", err)
		return c.Status(500).JSON(&fiber.Map{
			"error": "Internal Server Error",
		})
	}

	return c.Status(200).JSON(&fiber.Map{
		"count": len(resp.Users),
	})
}
