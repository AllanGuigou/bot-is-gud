package main

import (
	"context"
	"guigou/bot-is-gud/api/rpc"
	"net/http"
	"sort"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

type Presence struct {
	logger    *zap.SugaredLogger
	db        *pgxpool.Pool
	ctx       context.Context
	cache     *cache.Cache
	client    rpc.Presence
	isActive  bool
	startedAt time.Time
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

	p.startedAt = time.Now().UTC()
	p.logger.Info("presence service ready")
	return p
}

func (p *Presence) Stats() time.Duration {
	if !p.isActive {
		return 0
	}

	return time.Now().UTC().Sub(p.startedAt)
}

func (p *Presence) IsHealthy() bool {
	return p.isActive
}

type User struct {
	UID         string
	HasPresence bool
	Duration    time.Duration
}

func percentageToRank(percentage float64) string {
	switch {
	case percentage > .75:
		return "1"
	case percentage > .5:
		return "2"
	default:
		return "3"
	}
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
			UID:         uid,
			HasPresence: false,
			Duration:    time.Duration(0),
		}
	}

	return &User{
		UID:         uid,
		HasPresence: true,
		// add 30 seconds so that we can round up to the nearest minute and avoid zero durations for valid users
		Duration: time.Now().UTC().Add(time.Second * 30).Sub(start).Round(time.Minute),
	}
}

func (p *Presence) GetRecentUsers() []*User {
	after := time.Now().UTC().Add(-24 * time.Hour)
	return p.getUsers(after)
}

func (p *Presence) GetLeaderUsers(after time.Time) map[string][]string {
	if after.After(time.Now().UTC()) {
		return nil
	}

	users := p.getUsers(after)

	// get the duration of the top user
	var topDuration time.Duration
	if len(users) > 0 {
		tu := users[1]
		topDuration = tu.Duration
	}

	leaders := make(map[string][]string, 0)
	for _, user := range users {
		percentage := float64(user.Duration) / float64(topDuration)
		rank := percentageToRank(percentage)
		leaders[rank] = append(leaders[rank], user.UID)
	}

	return leaders
}

func (p *Presence) getUsers(after time.Time) []*User {
	users := make([]*User, 0)
	// presences could have an expire after `after` but would be excluded if their start was before `after`
	rows, err := p.db.Query(p.ctx, `
		SELECT uid, SUM(EXTRACT(EPOCH FROM (expire - start))) as duration, MAX(expire) as last_expiry 	
		FROM presences
		WHERE start > $1
		GROUP BY uid
	`, after)
	if err != nil {
		p.logger.Error(err)
		return users
	}

	um := make(map[string]*User)
	for rows.Next() {
		var uid string
		var td float64
		var le time.Time
		err = rows.Scan(&uid, &td, &le)
		if err != nil {
			p.logger.Error(err)
			continue
		}

		user, ok := um[uid]
		if !ok {
			user = &User{
				UID:         uid,
				HasPresence: false,
				Duration:    time.Duration(0),
			}
			um[uid] = user
		}

		user.Duration = time.Duration(td) * time.Second
		after := time.Now().UTC().Add(-3 * time.Minute)
		if le.After(after) {
			user.HasPresence = true
		}
	}

	for _, u := range um {
		users = append(users, u)
	}

	sort.Slice(users, func(i int, j int) bool {
		return users[i].Duration > users[j].Duration
	})

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
