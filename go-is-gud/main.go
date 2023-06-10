package main

import (
	"context"
	"fmt"
	"guigou/bot-is-gud/api"
	"guigou/bot-is-gud/api/rpc"
	"guigou/bot-is-gud/db"
	"guigou/bot-is-gud/env"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	env.Parse()
}

type Event struct {
	timestamp time.Time
	user      string
	action    string
}

var slash *Slash

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("go-is-gud is starting up...")

	dg, err := discordgo.New("Bot " + env.Token)
	if err != nil {
		fmt.Println(err)
		return
	}

	dg.Identify.Intents =
		discordgo.IntentsMessageContent +
			discordgo.IntentsDirectMessages +
			discordgo.IntentsGuilds +
			discordgo.IntentsGuildMessages +
			discordgo.IntentsGuildMessageTyping +
			discordgo.IntentsGuildVoiceStates +
			discordgo.IntentsGuildMembers

	err = dg.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("discordgo service is ready")

	lastTypedAt := NewTyper(dg)

	ctx := context.Background()
	api := api.New(env.PORT, lastTypedAt, ctx)

	db := db.New(ctx)
	rpc.SetupPresenceServer(dg, env.GID)
	if db != nil {
		p := New(ctx, db)
		api.RegisterHealthCheck(func() bool { return p.IsHealthy() })
		slash = NewSlash(dg, p)
		NewSuper(dg, p, env.SUID)
	} else {
		slash = NewSlash(dg, nil)
		NewSuper(dg, nil, env.SUID)
	}

	if env.ENABLE_BIGLY {
		dg.ApplicationCommandCreate(dg.State.User.ID, "", &discordgo.ApplicationCommand{
			Name:        "profile",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Configure a profile.",
		})
	}

	fmt.Println("go-is-gud is ready")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
	fmt.Println()
}
