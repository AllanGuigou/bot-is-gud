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
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	env.Parse()
}

var LastTypedAt time.Time = time.Unix(0, 0) // not thread safe but no big deal if this triggers twice

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

	dg.AddHandler(typingStart)
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

	ctx := context.Background()
	api := api.New(env.PORT, &LastTypedAt, ctx)

	db := db.New(ctx)
	rpc.SetupPresenceServer(dg, env.GID)
	if db != nil {
		p := New(ctx, db)
		api.RegisterHealthCheck(func() bool { return p.IsHealthy() })
		dg.AddHandler(messageCreate(p))
		slash = NewSlash(dg, p)
	} else {
		dg.AddHandler(messageCreate(nil))
		slash = NewSlash(dg, nil)
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

func triggerTyping(s *discordgo.Session, cid string) {
	if LastTypedAt.Add(time.Minute).After(time.Now().UTC()) {
		fmt.Println("typing too soon")
		return
	}

	LastTypedAt = time.Now().UTC()
	if rand.Intn(100) > 20 {
		fmt.Println("typing skipped")
		return
	}
	err := s.ChannelTyping(cid)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("typing triggered")
}

func typingStart(s *discordgo.Session, m *discordgo.TypingStart) {
	// ignore all messages created by the bot itself
	if m.UserID == s.State.User.ID {
		return
	}

	triggerTyping(s, m.ChannelID)
}

func messageCreate(p *Presence) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// ignore all messages created by the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Author.ID == env.SUID {
			guildID := m.GuildID
			if strings.HasPrefix(m.Content, ".disable") || strings.HasPrefix(m.Content, ".enable") {
				contents := strings.SplitAfter(m.Content, " ")
				if len(contents) < 2 {
					return
				}
				action := strings.TrimSpace(contents[0])
				command := strings.TrimSpace(contents[1])
				description := "no description"
				if len(contents) > 2 {
					description = strings.Join(contents[2:], " ")
				}
				switch action {
				case ".disable":
					{
						slash.remove(command, guildID)
					}
				case ".enable":
					{
						slash.add(command, description, guildID, make([]*discordgo.ApplicationCommandOption, 0))
					}
				}
			}

			if strings.HasPrefix(m.Content, ".user") {
				contents := strings.SplitAfter(m.Content, " ")
				if len(contents) < 2 {
					return
				}
				uid := contents[1]
				user, err := s.User(uid)
				if err != nil {
					fmt.Println(err)
					s.ChannelMessageSendReply(m.ChannelID, "error finding user", m.Reference())
					return
				}

				presence := p.GetUser(uid)
				if presence != nil && presence.HasPresence {
					s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("found user %s active for %s", user, presence.Duration), m.Reference())
				} else {
					s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("found user %s", user), m.Reference())
				}
			}
		}

		triggerTyping(s, m.ChannelID)
	}
}
