package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Super struct {
	dg   *discordgo.Session
	p    *Presence
	suid string
}

func NewSuper(dg *discordgo.Session, p *Presence, suid string) *Super {
	s := &Super{dg: dg, p: p, suid: suid}
	if suid == "" {
		fmt.Println("super user commands not available")
		return nil
	}
	dg.AddHandler(s.messageCreate)
	return s
}

func (s *Super) messageCreate(dg *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all messages created by non super users or the bot itself
	if m.Author.ID == dg.State.User.ID || m.Author.ID != s.suid {
		return
	}

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
		user, err := dg.User(uid)
		if err != nil {
			fmt.Println(err)
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
}
