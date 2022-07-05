package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type Slash struct {
	dg       *discordgo.Session
	mu       sync.Mutex
	commands map[string]*discordgo.ApplicationCommand
}

func NewSlash(dg *discordgo.Session) *Slash {
	// TODO: how to keep state of all commands at initialization
	s := Slash{
		dg:       dg,
		mu:       sync.Mutex{},
		commands: make(map[string]*discordgo.ApplicationCommand),
	}

	return &s
}

func (s *Slash) add(name, description, guildId string, options []*discordgo.ApplicationCommandOption) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ac := &discordgo.ApplicationCommand{
		Name:        name,
		Type:        discordgo.ChatApplicationCommand,
		Description: description,
		Options:     options,
	}

	_, err := s.dg.ApplicationCommandCreate(s.dg.State.User.ID, guildId, ac)
	if err != nil {
		panic(err)
	}
	s.commands[name] = ac

	log.Println(fmt.Sprintf("Added %s slash command to guild '%s'", name, guildId))
}

func (s *Slash) remove(name, guildId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	commands := make([]*discordgo.ApplicationCommand, 0)
	for k, v := range s.commands {
		if k != name {
			commands = append(commands, v)
		}
	}
	// bulk overwrite seems to cause commands to be removed immediately
	_, err := s.dg.ApplicationCommandBulkOverwrite(s.dg.State.User.ID, guildId, commands)
	if err != nil {
		panic(err)
	}

	delete(s.commands, name)
	log.Println(fmt.Sprintf("Removed %s slash command from guild '%s'", name, guildId))
}
