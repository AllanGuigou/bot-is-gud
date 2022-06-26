package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Profile struct {
	userId   string
	birthday time.Time
}

func NewProfile(d discordgo.ModalSubmitInteractionData) (*Profile, error) {
	if !strings.HasPrefix(d.CustomID, "profile_") {
		return nil, errors.New(fmt.Sprintf("Unable to parse profile from ModalSubmitInteractionData with CustomId:%s", d.CustomID))
	}

	// TODO: grab userId from interaction
	u := strings.Split(d.CustomID, "_")[1]

	month, err := parseTextInputAsNumber(d.Components[0])

	if err != nil {
		return nil, err
	}

	day, err := parseTextInputAsNumber(d.Components[1])

	if err != nil {
		return nil, err
	}

	p := Profile{userId: u, birthday: time.Date(1970, time.Month(month), day, 0, 0, 0, 0, time.UTC)}
	return &p, nil
}

func (p Profile) String() string {
	return fmt.Sprintf("%s has a birthday on %s", p.userId, p.birthday)
}

func parseTextInputAsNumber(r discordgo.MessageComponent) (int, error) {
	v := r.(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	return strconv.Atoi(v)
}
