package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

var lastTypedAt time.Time = time.Unix(0, 0) // not thread safe but no big deal if this triggers twice

func NewTyper(dg *discordgo.Session) *time.Time {
	dg.AddHandler(typingStart)
	dg.AddHandler(messageCreate)
	return &lastTypedAt
}

func triggerTyping(s *discordgo.Session, cid string) {
	if lastTypedAt.Add(time.Minute).After(time.Now().UTC()) {
		fmt.Println("typing too soon")
		return
	}

	lastTypedAt = time.Now().UTC()
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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	triggerTyping(s, m.ChannelID)
}
