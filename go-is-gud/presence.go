package main

import (
	"context"
	"guigou/bot-is-gud/api/rpc"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

type Presence struct {
	logger   *zap.SugaredLogger
	db       *pgxpool.Pool
	ctx      context.Context
	cache    *cache.Cache
	client   rpc.Presence
	isActive bool
}

func New(logger *zap.SugaredLogger, ctx context.Context, db *pgxpool.Pool) *Presence {
	client := rpc.NewPresenceProtobufClient("http://localhost:8080", &http.Client{})
	c := cache.New(2*time.Minute, 2*time.Minute)
	c.OnEvicted(traceInactive(logger))
	p := &Presence{logger: logger, db: db, ctx: ctx, cache: c, client: client, isActive: true}

	go func() {
		p.record()
		p.logger.Error("presence service is inactive")
		p.isActive = false
	}()

	p.logger.Info("presence service ready")
	return p
}

func (p *Presence) IsHealthy() bool {
	return p.isActive
}

type User struct {
	HasPresence bool
	Duration    time.Duration
}

func (p *Presence) GetUser(uid string) *User {
	if p == nil {
		return nil
	}

	var start time.Time
	after := time.Now().UTC().Add(-3 * time.Minute)
	err := p.db.QueryRow(p.ctx, `SELECT start FROM presences WHERE expire > $1 AND uid = $2`, after, uid).Scan(&start)
	if err != nil {
		if err != pgx.ErrNoRows {
			p.logger.Error(err)
			return nil
		}
		return &User{
			HasPresence: false,
			Duration:    time.Duration(0),
		}
	}

	return &User{
		HasPresence: true,
		// add 30 seconds so that we can round up to the nearest minute and avoid zero durations for valid users
		Duration: time.Now().UTC().Add(time.Second * 30).Sub(start).Round(time.Minute),
	}
}

func (p *Presence) GetRecentUsers() map[string]*User {
	users := make(map[string]*User)
	after := time.Now().UTC().Add(-24 * time.Hour)
	// presences could have an expire after `after` but would be excluded if their start was before `after`
	rows, err := p.db.Query(p.ctx, `SELECT uid, start, expire FROM presences WHERE start > $1`, after)
	if err != nil {
		p.logger.Error(err)
		return users
	}

	for rows.Next() {
		var uid string
		var start time.Time
		var expire time.Time
		err = rows.Scan(&uid, &start, &expire)
		if err != nil {
			p.logger.Error(err)
		}

		user, ok := users[uid]
		if !ok {
			user = &User{
				HasPresence: false,
				Duration:    time.Duration(0),
			}
			users[uid] = user
		}

		user.Duration += expire.Sub(start)
		after := time.Now().UTC().Add(-3 * time.Minute)
		if expire.After(after) {
			user.HasPresence = true
		}
	}

	return users
}

func traceInactive(logger *zap.SugaredLogger) func(uid string, id interface{}) {
	return func(uid string, id interface{}) {
		logger.Infow("removing active presence", "uid", uid)
	}
}

func (p *Presence) record() {
	after := time.Now().UTC().Add(-3 * time.Minute)
	rows, err := p.db.Query(p.ctx, `SELECT id, uid, expire FROM presences WHERE expire > $1`, after)
	if err != nil {
		p.logger.Error(err)
		return
	}

	for rows.Next() {
		var id int64
		var uid string
		var expire time.Time
		err = rows.Scan(&id, &uid, &expire)
		if err != nil {
			p.logger.Error(err)
			return
		}

		p.cache.Set(uid, id, 0)
	}

	for {
		now := time.Now().UTC()
		res, err := p.client.WhoseOn(p.ctx, &rpc.WhoseOnReq{VoiceChannel: ""})
		if err != nil {
			p.logger.Error(err)
			return
		}

		// TODO: update to use copy from to bulk insert
		// inserts := make([][]interface{}, 0)
		updates := make([]int64, 0)

		for _, u := range res.Users {
			id, found := p.cache.Get(u)
			if found {
				updates = append(updates, id.(int64))
			} else {
				p.logger.Infow("creating active presence", "uid", u)
				// TODO: deprecate `active`
				err := p.db.QueryRow(p.ctx, `INSERT INTO presences (uid, active, start, expire) VALUES ($1, TRUE, $2, $3) RETURNING (id)`, u, now, now).Scan(&id)
				if err != nil {
					p.logger.Error(err)
					return
				}
			}
			// set cache for new users and refresh for existing users
			p.cache.Set(u, id, 0)
		}

		if len(updates) > 0 {
			p.logger.Infow("updating active presence", zap.Int64s("uids", updates))
			_, err := p.db.Exec(p.ctx, `UPDATE presences SET expire = $1 WHERE id = ANY ($2)`, now, updates)
			if err != nil {
				p.logger.Error(err)
				return
			}
		}

		time.Sleep(1 * time.Minute)
	}
}
