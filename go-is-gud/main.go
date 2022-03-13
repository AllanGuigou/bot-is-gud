package main

import (
	"flag"
	"fmt"
	"guigou/bot-is-gud/health"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", LookupEnvOrString("DISCORD_TOKEN", Token), "Bot Token")
	flag.Parse()
}

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return defaultVal
}

// not thread safe but no big deal if this triggers twice
var LastTypedAt time.Time = time.Unix(0, 0)

type Event struct {
	timestamp time.Time
	user      string
	action    string
}

func track(c <-chan Event) {
	for event := range c {
		fmt.Printf("%s %s %s\n", event.timestamp, event.user, event.action)
	}
}

func main() {
	go health.New(&LastTypedAt)
	c := make(chan Event, 100)
	go track(c)

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println(err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(typingStart)
	dg.AddHandler(slashCommandHandler(c))
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	dg.Identify.Intents = discordgo.IntentsGuildMessageTyping

	command := &discordgo.ApplicationCommand{
		Name:        "bigly",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Word of the day!",
	}

	err = dg.Open()
	if err != nil {
		fmt.Println(err)
		return
	}

	dg.ApplicationCommandCreate(dg.State.User.ID, "", command)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
	fmt.Println()
}

// rand.Seed(time.Now().UnixNano())
var wl []string = []string{
	"aback",
	"abase",
	"abate",
	"abbey",
	"abbot",
	"abhor",
	"abide",
	"abled",
	"abode",
	"abort",
}

func rw() string {
	i := rand.Intn(len(wl))
	return wl[i]
}

func eventFromInteraction(i discordgo.InteractionCreate) Event {
	var user string
	if i.Member != nil {
		user = i.Member.User.ID
	}

	user = i.User.ID

	// TODO: how to get timestamp from interaction
	return Event{timestamp: time.Now(), user: user, action: "bigly-slash-command"}
}

func slashCommandHandler(c chan<- Event) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		command := i.ApplicationCommandData()
		switch cn := command.Name; cn {
		case "bigly":
			{
				content := rw()
				c <- eventFromInteraction(*i)
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: content,
						Flags:   uint64(discordgo.MessageFlagsEphemeral),
					},
				})

				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func triggerTyping(s *discordgo.Session, cid string) {
	if LastTypedAt.Add(time.Minute).After(time.Now().UTC()) {
		fmt.Println("typing too soon")
		return
	}

	LastTypedAt = time.Now().UTC()
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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	triggerTyping(s, m.ChannelID)
}
