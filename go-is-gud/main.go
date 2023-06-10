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
	"go.uber.org/ratelimit"
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

func track(c <-chan Event) {
	for event := range c {
		fmt.Printf("%s %s %s\n", event.timestamp, event.user, event.action)
	}
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
	c := make(chan Event, 100)
	go track(c)
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
		dg.AddHandler(slashCommandHandler(c, p))
	} else {
		dg.AddHandler(messageCreate(nil))
		dg.AddHandler(slashCommandHandler(c, nil))
	}

	slash = NewSlash(dg)

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

func slashCommandHandler(c chan<- Event, p *Presence) func(*discordgo.Session, *discordgo.InteractionCreate) {
	rl := ratelimit.New(5, ratelimit.Per(1*time.Minute))
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			{
				command := i.ApplicationCommandData()
				switch cn := command.Name; cn {
				case "me":
					{
						event := eventFromInteraction(*i)
						c <- event
						err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Flags: uint64(discordgo.MessageFlagsEphemeral),
							},
						})

						if err != nil {
							fmt.Println(err)
							return
						}

						rl.Take()
						var content string
						// TODO: use dg state to provide more accurate results if the user recently connected or disconnected
						presence := p.GetUser(event.user)
						if presence == nil {
							content = "Oops something went wrong..."
							_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
								Content: content,
							})

							if err != nil {
								fmt.Println(err)
								return
							}
						} else if presence.HasPresence {
							content = fmt.Sprintf("You've been active for %s", strings.ReplaceAll(presence.Duration.String(), "0s", ""))
							_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
								Content: content,
							})

							if err != nil {
								fmt.Println(err)
								return
							}
						} else {
							content = "You've been inactive"
							_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
								Content: content,
							})

							if err != nil {
								fmt.Println(err)
								return
							}
						}
						fmt.Printf("'/me' slash command used for user: '%s' content: '%s'\n", event.user, content)
					}
				case "profile":
					{
						c <- eventFromInteraction(*i)
						err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseModal,
							Data: &discordgo.InteractionResponseData{
								CustomID: "profile_" + i.Interaction.Member.User.ID,
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
				case "lets-gamble":
					{
						err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
