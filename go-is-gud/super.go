package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type Super struct {
	logger   *zap.SugaredLogger
	dg       *discordgo.Session
	p        *Presence
	suid     string
	commands []string
}

func NewSuper(logger *zap.SugaredLogger, dg *discordgo.Session, p *Presence, suid string) *Super {
	s := &Super{logger: logger, dg: dg, p: p, suid: suid, commands: []string{".enable", ".disable", ".user", ".users", ".restart"}}
	if suid == "" {
		logger.Warn("failed to setup super commands")
		return nil
	}
	dg.AddHandler(s.messageCreate)
	return s
}

func Contains(a []string, s string) bool {
	for _, c := range a {
		if c == s {
			return true
		}
	}
	return false
}

func (s *Super) messageCreate(dg *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all messages created by non super users or the bot itself
	if m.Author.ID == dg.State.User.ID || m.Author.ID != s.suid {
		return
	}

	contents := strings.SplitAfter(m.Content, " ")
	if len(contents) == 0 || !Contains(s.commands, strings.TrimSpace(contents[0])) {
		return
	}

	command := strings.TrimSpace(contents[0])
	s.logger.Infow("super command triggered",
		"name", command)

	switch command {
	case ".disable":
		{
			if len(contents) < 2 {
				return
			}
			command := strings.TrimSpace(contents[1])
			slash.remove(command, m.GuildID)
		}
	case ".enable":
		{
			if len(contents) < 2 {
				return
			}
			command := strings.TrimSpace(contents[1])
			description := "no description"
			if len(contents) > 2 {
				description = strings.Join(contents[2:], " ")
			}
			slash.add(command, description, m.GuildID, make([]*discordgo.ApplicationCommandOption, 0))
		}
	case ".user":
		{
			if len(contents) < 2 {
				return
			}
			uid := contents[1]
			user, err := dg.User(uid)
			if err != nil {
				s.logger.Error(err)
				dg.ChannelMessageSendReply(m.ChannelID, "error finding user", m.Reference())
				return
			}
			presence := s.p.GetUser(uid)
			if presence != nil && presence.HasPresence {
				dg.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("found user %s active for %s", user, presence.Duration), m.Reference())
			} else {
				dg.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("found user %s", user), m.Reference())
			}
		}
	case ".users":
		{
			var duration time.Duration = 24 * time.Hour
			if len(contents) > 1 {
				d, err := time.ParseDuration(contents[1])
				if err != nil {
					dg.ChannelMessageSendReply(m.ChannelID, "invalid duration format", m.Reference())
					return
				}
				duration = d.Abs()
			}

			after := time.Now().UTC().Add(-duration)
			leaders := s.p.GetLeaderUsers(after)
			if len(leaders) < 1 {
				dg.ChannelMessageSendReply(m.ChannelID, "no users found", m.Reference())
			} else {
				message := ""
				ranks := make([]string, 0, len(leaders))
				for rank := range leaders {
					ranks = append(ranks, rank)
				}

				for _, r := range ranks {
					message += fmt.Sprintf("%s:", r)
					for _, u := range leaders[r] {
						username := getUsername(dg, u)
						message += fmt.Sprintf(" %s", username)
					}
					message += fmt.Sprintln()
					// status := ""
					// if u.HasPresence {
					// 	status = "(active)"
					// }
					// message += fmt.Sprintf("%s %s %s\n", username, u.Duration, status)
				}
				dg.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
			}
		}
	case ".restart":
		{
			p, err := os.FindProcess(os.Getpid())
			if err != nil {
				s.logger.Error(err)
				return
			}
			p.Signal(syscall.SIGINT)
		}
	}
}

func getUsername(dg *discordgo.Session, uid string) string {
	user, err := dg.User(uid)
	if err != nil {
		return fmt.Sprintf("<unknown-%s>", uid)
	}
	return user.GlobalName
}
