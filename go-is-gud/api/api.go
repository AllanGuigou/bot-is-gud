package api

import (
	"context"
	"guigou/bot-is-gud/api/rpc"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type API struct {
	lastTypedAt *time.Time
	protoClient rpc.Presence
}

func New(lastTypedAt *time.Time, port string) {
	app := fiber.New()
	client := rpc.NewPresenceProtobufClient("http://localhost:8080", &http.Client{})
	api := &API{lastTypedAt: lastTypedAt, protoClient: client}
	app.Get("/", api.healthHandler)
	app.Get("/whoseOn", api.whoseOnHandler)
	app.Listen(":" + port)
}

func (api *API) healthHandler(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(200).JSON(&fiber.Map{
		"lastTypedAt": api.lastTypedAt.Format(time.RFC3339),
		"status":      "Ok",
	})
}

func (api *API) whoseOnHandler(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	resp, err := api.protoClient.WhoseOn(context.Background(), &rpc.WhoseOnReq{VoiceChannel: ""})
	if err != nil {
		return err
	}

	return c.Status(200).JSON(&fiber.Map{
		"count": len(resp.Users),
	})
}
