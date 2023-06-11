package main

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var lastTypedAt time.Time = time.Unix(0, 0) // not thread safe but no big deal if this triggers twice

func NewTyper(logger *zap.SugaredLogger, dg *discordgo.Session) *time.Time {
	dg.AddHandler(typingStart(logger))
	dg.AddHandler(messageCreate(logger))
	return &lastTypedAt
}

func triggerTyping(logger *zap.SugaredLogger, s *discordgo.Session, cid string) {
	if lastTypedAt.Add(time.Minute).After(time.Now().UTC()) {
		logger.Debug("typing too soon")
		return
	}

	lastTypedAt = time.Now().UTC()
	if rand.Intn(100) > 20 {
		logger.Debug("typing skipped")
		return
	}
	err := s.ChannelTyping(cid)

	if err != nil {
		logger.Error(err)
		return
	}

	logger.Infow("typing triggered",
		"cid", cid)
}

func typingStart(logger *zap.SugaredLogger) func(s *discordgo.Session, m *discordgo.TypingStart) {
	return func(s *discordgo.Session, m *discordgo.TypingStart) {
		// ignore all messages created by the bot itself
		if m.UserID == s.State.User.ID {
			return
		}

		triggerTyping(logger, s, m.ChannelID)
	}
}

func messageCreate(logger *zap.SugaredLogger) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// ignore all messages created by the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		triggerTyping(logger, s, m.ChannelID)
	}
}
