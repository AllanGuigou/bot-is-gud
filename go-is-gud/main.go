package main

import (
	"context"
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
	"go.uber.org/zap"
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
	// TODO: consider adding logger to ctx to avoid having to pass it around individually
	logger := NewLogger()
	defer logger.Sync()
	logger.Info("go-is-gud is starting up...")

	dg, err := setupdg(logger)
	if err != nil {
		logger.Fatal(err)
	}

	// db
	ctx := context.Background()
	db := db.New(logger, ctx)

	// features
	lastTypedAt := NewTyper(logger, dg)

	api := api.New(logger, env.PORT, lastTypedAt, ctx)

	rpc.SetupPresenceServer(logger, dg, env.GID)
	if db != nil {
		p := New(logger, ctx, db)
		api.RegisterHealthCheck(func() bool { return p.IsHealthy() })
		slash = NewSlash(logger, dg, p)
		NewSuper(logger, dg, p, env.SUID)
	} else {
		slash = NewSlash(logger, dg, nil)
		NewSuper(logger, dg, nil, env.SUID)
	}

	if env.ENABLE_BIGLY {
		dg.ApplicationCommandCreate(dg.State.User.ID, "", &discordgo.ApplicationCommand{
			Name:        "profile",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Configure a profile.",
		})
	}

	logger.Info("go-is-gud is ready")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func setupdg(logger *zap.SugaredLogger) (*discordgo.Session, error) {
	dg, err := discordgo.New("Bot " + env.Token)
	if err != nil {
		logger.Info("failed to setup discordgo")
		return nil, err
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
		logger.Info("failed to setup discordgo")
		return nil, err
	}

	logger.Info("discordgo service is ready")
	return dg, nil
}
