package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
)

type Slash struct {
	logger   *zap.SugaredLogger
	dg       *discordgo.Session
	mu       sync.Mutex
	commands map[string]*discordgo.ApplicationCommand

	presence *Presence
}

func NewSlash(logger *zap.SugaredLogger, dg *discordgo.Session, p *Presence) *Slash {
	// TODO: how to keep state of all commands at initialization
	s := Slash{
		logger:   logger,
		dg:       dg,
		mu:       sync.Mutex{},
		commands: make(map[string]*discordgo.ApplicationCommand),
		presence: p,
	}

	dg.AddHandler(s.commandHandler)

	return &s
}

func (s *Slash) add(name, description, guildID string, options []*discordgo.ApplicationCommandOption) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ac := &discordgo.ApplicationCommand{
		Name:        name,
		Type:        discordgo.ChatApplicationCommand,
		Description: description,
		Options:     options,
	}

	s.logger.Infow("Adding slash command",
		"gid", guildID,
		"name", name,
		"description", description)

	_, err := s.dg.ApplicationCommandCreate(s.dg.State.User.ID, guildID, ac)
	if err != nil {
		s.logger.Panic(err)
	}
	s.commands[name] = ac
}

func (s *Slash) remove(name, guildID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	commands := make([]*discordgo.ApplicationCommand, 0)
	for k, v := range s.commands {
		if k != name {
			commands = append(commands, v)
		}
	}

	s.logger.Infow("Removing slash command",
		"gid", guildID,
		"name", name)

	// bulk overwrite seems to cause commands to be removed immediately
	_, err := s.dg.ApplicationCommandBulkOverwrite(s.dg.State.User.ID, guildID, commands)
	if err != nil {
		panic(err)
	}

	delete(s.commands, name)
}

func eventFromInteraction(i discordgo.InteractionCreate) Event {
	var user string
	if i.Member != nil {
		user = i.Member.User.ID
	} else if i.User != nil {
		user = i.User.ID
	}

	// TODO: how to get timestamp from interaction
	return Event{timestamp: time.Now(), user: user, action: i.ApplicationCommandData().Name}
}

func (slash *Slash) commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ratelimitMe := ratelimit.New(5, ratelimit.Per(1*time.Minute))
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		{
			command := i.ApplicationCommandData()
			event := eventFromInteraction(*i)
			slash.logger.Infow("slash command triggered", "name", event.action, "uid", event.user)
			switch cn := command.Name; cn {
			case "me":
				{
					meCommand(slash.logger, s, i.Interaction, event, ratelimitMe, slash.presence)
				}
			case "stats":
				{
					statsCommand(slash.logger, s, i.Interaction, ratelimitMe, slash.presence)
				}
			case "profile":
				{
					profileCommand(slash.logger, s, i.Interaction)
				}
			case "lets-gamble":
				{
					gambleCommand(slash.logger, s, i.Interaction)
				}
			}
		}
	case discordgo.InteractionModalSubmit:
		{
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: 1 << 6,
				},
			})

			if err != nil {
				slash.logger.Error(err)
			}

			p, err := NewProfile(i.ModalSubmitData())
			if err != nil {
				slash.logger.Error(err)
				return
			}

			slash.logger.Info(p)
		}
	}
}

func meCommand(logger *zap.SugaredLogger, s *discordgo.Session, i *discordgo.Interaction, event Event, ratelimit ratelimit.Limiter, presence *Presence) {
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		logger.Error(err)
		return
	}

	ratelimit.Take()
	var content string
	// TODO: use dg state to provide more accurate results if the user recently connected or disconnected
	up := presence.GetUser(event.user)
	if up == nil {
		content = "Oops something went wrong..."
		_, err = s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
			Content: content,
		})

		if err != nil {
			logger.Error(err)
			return
		}
	} else if up.HasPresence {
		content = fmt.Sprintf("You've been active for %s", strings.ReplaceAll(up.Duration.String(), "0s", ""))
		_, err = s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
			Content: content,
		})

		if err != nil {
			logger.Error(err)
			return
		}
	} else {
		content = "You've been inactive"
		_, err = s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
			Content: content,
		})

		if err != nil {
			logger.Error(err)
			return
		}
	}

	logger.Infow("slash command completed",
		"name", "me",
		"uid", event.user,
		"content", content)
}

func statsCommand(logger *zap.SugaredLogger, s *discordgo.Session, i *discordgo.Interaction, ratelimit ratelimit.Limiter, presence *Presence) {
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		logger.Error(err)
		return
	}

	ratelimit.Take()
	var content string
	if presence == nil {
		content = "```\nUptime: unavailable\n```"
	} else {
		uptime := presence.Stats()
		content = fmt.Sprintf("```\nUptime: %s\n```", uptime)
	}
	_, err = s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
		Content: content,
	})

	if err != nil {
		logger.Error(err)
		return
	}

	logger.Infow("slash command completed",
		"name", "status",
		"content", content)
}

func profileCommand(logger *zap.SugaredLogger, s *discordgo.Session, i *discordgo.Interaction) {
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "profile_" + i.Member.User.ID,
			Title:    "Profile",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "birthday_month",
							Label:    "Month",
							Style:    discordgo.TextInputShort,
							Required: true,
							Value:    "07",
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "birthday_day",
							Label:    "Day",
							Style:    discordgo.TextInputShort,
							Required: true,
							Value:    "24",
						},
					},
				},
			},
		}},
	)

	if err != nil {
		logger.Error(err)
	}
}

// gambleLosses tracks consecutive global losses to increase win odds.
var gambleLosses = struct {
	sync.Mutex
	count int
}{}

func gambleCommand(logger *zap.SugaredLogger, s *discordgo.Session, i *discordgo.Interaction) {
	if i.Member == nil || i.Member.User == nil {
		logger.Info("attempt to `lets-gamble` from outside of a guild")
		return
	}

	if i.Member.CommunicationDisabledUntil != nil && i.Member.CommunicationDisabledUntil.After(time.Now()) {
		s.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You're already timed out. Nice try.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	uid := i.Member.User.ID

	// Base win chance 5%, +5% per consecutive global loss, capped at 30%
	gambleLosses.Lock()
	losses := gambleLosses.count
	winChance := 5 + (losses * 5)
	if winChance > 30 {
		winChance = 30
	}
	won := rand.Intn(100) < winChance
	if won {
		gambleLosses.count = 0
	} else {
		gambleLosses.count++
	}
	gambleLosses.Unlock()

	var content string
	if won {
		logger.Infow("gamble won", "uid", uid, "streak", losses)
		content = "You win! 🎉 ...this time."
	} else {
		t := time.Now().Add(time.Second * time.Duration(rand.Intn(1000)))
		err := s.GuildMemberTimeout(i.GuildID, uid, &t)
		if err != nil {
			logger.Infow("gamble won (timeout failed)", "uid", uid, "error", err)
			gambleLosses.Lock()
			gambleLosses.count = 0
			gambleLosses.Unlock()
			content = "You win! 🎉 ...this time."
		} else {
			content = "You've lost."
		}
	}

	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		logger.Error(err)
	}
}
