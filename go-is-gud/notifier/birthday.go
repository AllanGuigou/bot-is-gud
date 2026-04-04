package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Birthday struct {
	logger  *zap.SugaredLogger
	dg      *discordgo.Session
	db      *pgxpool.Pool
	ctx     context.Context
	suid    string
	trigger chan struct{}
}

func NewBirthday(logger *zap.SugaredLogger, dg *discordgo.Session, db *pgxpool.Pool, suid string, ctx context.Context) *Birthday {
	b := &Birthday{
		logger:  logger,
		dg:      dg,
		db:      db,
		ctx:     ctx,
		suid:    suid,
		trigger: make(chan struct{}, 1),
	}

	go Schedule(ctx, 24*time.Hour, 22*time.Hour, b.trigger, b.check)

	logger.Info("birthday notifier ready")
	return b
}

// Trigger signals the schedule to immediately process a birthday check.
func (b *Birthday) Trigger() {
	select {
	case b.trigger <- struct{}{}:
	default:
	}
}

func (b *Birthday) check(_ time.Time) {
	if b.db == nil {
		b.logger.Warn("birthday notifier has no database connection")
		return
	}

	rows, err := b.db.Query(b.ctx, `
		SELECT user_id FROM profiles
		WHERE EXTRACT(MONTH FROM birthday) = EXTRACT(MONTH FROM now() AT TIME ZONE timezone)
		  AND EXTRACT(DAY FROM birthday) = EXTRACT(DAY FROM now() AT TIME ZONE timezone)
	`)
	if err != nil {
		b.logger.Error(err)
		return
	}
	defer rows.Close()

	var sent []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			b.logger.Error(err)
			continue
		}

		channel, err := b.dg.UserChannelCreate(userID)
		if err != nil {
			b.logger.Errorw("failed to create DM channel", "uid", userID, "error", err)
			continue
		}

		_, err = b.dg.ChannelMessageSend(channel.ID, fmt.Sprintf("Happy Birthday <@%s> :birthday:", userID))
		if err != nil {
			b.logger.Errorw("failed to send birthday DM", "uid", userID, "error", err)
			continue
		}

		b.logger.Infow("birthday notification sent", "uid", userID)
		sent = append(sent, userID)
	}
	if err := rows.Err(); err != nil {
		b.logger.Error(err)
	}

	b.notifySuper(sent)
}

func (b *Birthday) notifySuper(sent []string) {
	if b.suid == "" {
		return
	}

	var message string
	if len(sent) == 0 {
		message = "birthday check completed, no birthdays today"
	} else {
		message = fmt.Sprintf("birthday notifications sent to %d user(s):", len(sent))
		for _, uid := range sent {
			message += fmt.Sprintf(" <@%s>", uid)
		}
	}

	channel, err := b.dg.UserChannelCreate(b.suid)
	if err != nil {
		b.logger.Errorw("failed to DM super user", "error", err)
		return
	}

	_, err = b.dg.ChannelMessageSend(channel.ID, message)
	if err != nil {
		b.logger.Errorw("failed to DM super user", "error", err)
	}
}
