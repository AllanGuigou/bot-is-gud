package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"

	"guigou/bot-is-gud/health"
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

func main() {
	health.New(&LastTypedAt)

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println(err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(typingStart)
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	dg.Identify.Intents = discordgo.IntentsGuildMessageTyping

	err = dg.Open()
	if err != nil {
		fmt.Println(err)
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
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
