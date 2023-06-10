package birthday

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v4"
)

type Notifier struct {
	dg  *discordgo.Session
	db  *pgx.Conn
	ctx context.Context
}

func New(dg *discordgo.Session, db *pgx.Conn, ctx context.Context) *Notifier {
	n := Notifier{dg: dg, db: db, ctx: ctx}

	// TODO: how to avoid multiple notifications or no notifications if the service restarts
	go Schedule(ctx, time.Hour*24, time.Hour*16, n.sendBirthdayMessage)

	return &n
}

func (n *Notifier) sendBirthdayMessage(time time.Time) {
	var userId string
	var channelId string

	// TODO: support multiple birthdays on a single day
	err := n.db.QueryRow(n.ctx, "SELECT userId, channelId FROM profiles WHERE date_trunc('month', birthday) = date_trunc('month', now()) AND date_trunc('day', birthday) = date_trunc('day', now())").Scan(&userId, &channelId)

	if err != nil && err != pgx.ErrNoRows {
		fmt.Println(err)
		return
	}

	// TODO: could potentially allow the user to create a server event for their birthday
	n.dg.ChannelMessageSend(channelId, fmt.Sprintf("Happy Birthday <@%s> :birthday:", userId))
}
