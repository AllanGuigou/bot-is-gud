package main

import (
	"fmt"
	"guigou/bot-is-gud/db"
	"guigou/bot-is-gud/env"
	"guigou/bot-is-gud/health"
	birthday "guigou/bot-is-gud/notifier"
	"log"
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

var slash *Slash

func main() {
	rand.Seed(time.Now().UnixNano())
	go health.New(&LastTypedAt, env.PORT)

	dg, err := discordgo.New("Bot " + env.Token)
	if err != nil {
		fmt.Println(err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(typingStart)
	c := make(chan Event, 100)
	go track(c)
	dg.AddHandler(slashCommandHandler(c))

	dg.Identify.Intents = discordgo.IntentsGuildMessages +
		discordgo.IntentsMessageContent +
		discordgo.IntentsDirectMessages +
		discordgo.IntentsGuildMessageTyping

	err = dg.Open()
	if err != nil {
		fmt.Println(err)
		return
	}

	db := db.New()
	go birthday.New(dg, db)
	slash = NewSlash(dg)

	if env.ENABLE_BIGLY {
		command := &discordgo.ApplicationCommand{
			Name:        "bigly",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Word of the day!",
		}
		dg.ApplicationCommandCreate(dg.State.User.ID, "", command)
		dg.ApplicationCommandCreate(dg.State.User.ID, "", &discordgo.ApplicationCommand{
			Name:        "profile",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Configure a profile.",
		})
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
	fmt.Println()
}

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
	} else if i.User != nil {
		user = i.User.ID
	}

	// TODO: how to get timestamp from interaction
	return Event{timestamp: time.Now(), user: user, action: "bigly-slash-command"}
}

func slashCommandHandler(c chan<- Event) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			{
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
					log.Println(err.Error())
					return
				}

				log.Println(p.String())
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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Author.ID == env.SUID {
		guildId := m.GuildID
		switch m.Content {
		case ".disable":
			slash.remove("lets-gamble", guildId)
		case ".enable":
			slash.add("lets-gamble", "...", guildId)
		}
	}

	triggerTyping(s, m.ChannelID)
}
