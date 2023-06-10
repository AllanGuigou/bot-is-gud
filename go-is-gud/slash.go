package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/ratelimit"
)

type Slash struct {
	dg       *discordgo.Session
	mu       sync.Mutex
	commands map[string]*discordgo.ApplicationCommand

	presence *Presence
}

func NewSlash(dg *discordgo.Session, p *Presence) *Slash {
	// TODO: how to keep state of all commands at initialization
	s := Slash{
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

	_, err := s.dg.ApplicationCommandCreate(s.dg.State.User.ID, guildID, ac)
	if err != nil {
		panic(err)
	}
	s.commands[name] = ac

	fmt.Printf("Added '%s' slash command to guild '%s' with description '%s'\n", name, guildID, description)
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
	// bulk overwrite seems to cause commands to be removed immediately
	_, err := s.dg.ApplicationCommandBulkOverwrite(s.dg.State.User.ID, guildID, commands)
	if err != nil {
		panic(err)
	}

	delete(s.commands, name)
	fmt.Printf("Removed %s slash command from guild '%s'\n", name, guildID)
}

func eventFromInteraction(i discordgo.InteractionCreate) Event {
	var user string
	if i.Member != nil {
		user = i.Member.User.ID
	} else if i.User != nil {
		user = i.User.ID
	}

	// TODO: how to get timestamp from interaction
	return Event{timestamp: time.Now(), user: user, action: "bigly-slash-command"}
}

func (slash *Slash) commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ratelimitMe := ratelimit.New(5, ratelimit.Per(1*time.Minute))
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		{
			command := i.ApplicationCommandData()
			event := eventFromInteraction(*i)
			fmt.Printf("%s %s %s\n", event.timestamp, event.user, event.action)
			switch cn := command.Name; cn {
			case "me":
				{
					meCommand(s, i.Interaction, event, ratelimitMe, slash.presence)
				}
			case "profile":
				{
					profileCommand(s, i.Interaction)
				}
			case "lets-gamble":
				{
					gambleCommand(s, i.Interaction)
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
				fmt.Println(err)
			}

			p, err := NewProfile(i.ModalSubmitData())
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(p)
		}
	}
}

func meCommand(s *discordgo.Session, i *discordgo.Interaction, event Event, ratelimit ratelimit.Limiter, presence *Presence) {
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: uint64(discordgo.MessageFlagsEphemeral),
		},
	})

	if err != nil {
		fmt.Println(err)
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
			fmt.Println(err)
			return
		}
	} else if up.HasPresence {
		content = fmt.Sprintf("You've been active for %s", strings.ReplaceAll(up.Duration.String(), "0s", ""))
		_, err = s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
			Content: content,
		})

		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		content = "You've been inactive"
		_, err = s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
			Content: content,
		})

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Printf("'/me' slash command used for user: '%s' content: '%s'\n", event.user, content)
}

func profileCommand(s *discordgo.Session, i *discordgo.Interaction) {
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
		fmt.Println(err)
	}
}

func gambleCommand(s *discordgo.Session, i *discordgo.Interaction) {
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You've lost.",
			Flags:   uint64(discordgo.MessageFlagsEphemeral),
		},
	},
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	t := time.Now().Add(time.Second * time.Duration(rand.Intn(1000)))

	if i.Member == nil || i.Member.User == nil {
		fmt.Println("Attempt to `lets-gamble` from outside of a guild.")
		return
	}

	err = s.GuildMemberTimeout(i.GuildID, i.Member.User.ID, &t)

	if err != nil {
		fmt.Println(err)
		return
	}
}
