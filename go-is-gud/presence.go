package main

import (
	"context"
	"fmt"
	"guigou/bot-is-gud/api/rpc"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickmn/go-cache"
)

type Presence struct {
	db       *pgxpool.Pool
	ctx      context.Context
	cache    *cache.Cache
	client   rpc.Presence
	isActive bool
}

func New(ctx context.Context, db *pgxpool.Pool) *Presence {
	client := rpc.NewPresenceProtobufClient("http://localhost:8080", &http.Client{})
	c := cache.New(2*time.Minute, 2*time.Minute)
	c.OnEvicted(traceInactive)
	p := &Presence{db: db, ctx: ctx, cache: c, client: client, isActive: true}

	go func() {
		p.record()
		p.isActive = false
	}()

	fmt.Println("presence service ready")
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
			fmt.Println(err)
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

func traceInactive(uid string, id interface{}) {
	fmt.Printf("removing %s active presence\n", uid)
}

func (p *Presence) record() {
	after := time.Now().UTC().Add(-3 * time.Minute)
	rows, err := p.db.Query(p.ctx, `SELECT id, uid, expire FROM presences WHERE expire > $1`, after)
	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		var id int64
		var uid string
		var expire time.Time
		err = rows.Scan(&id, &uid, &expire)
		if err != nil {
			fmt.Println(err)
			return
		}

		p.cache.Set(uid, id, 0)
	}

	for {
		now := time.Now().UTC()
		res, err := p.client.WhoseOn(p.ctx, &rpc.WhoseOnReq{VoiceChannel: ""})
		if err != nil {
			fmt.Println(err)
			return
		}

		// TODO: update to use copy from to bulk insert
		// inserts := make([][]interface{}, 0)
		updates := make([]int64, 0)

		for _, u := range res.Users {
			id, found := p.cache.Get(u)
			if found {
				fmt.Printf("updating %s active presence\n", u)
				updates = append(updates, id.(int64))
			} else {
				fmt.Printf("creating %s active presence\n", u)
				// TODO: deprecate `active`
				err := p.db.QueryRow(p.ctx, `INSERT INTO presences (uid, active, start, expire) VALUES ($1, TRUE, $2, $3) RETURNING (id)`, u, now, now).Scan(&id)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			// set cache for new users and refresh for existing users
			p.cache.Set(u, id, 0)
		}

		if len(updates) > 0 {
			_, err := p.db.Exec(p.ctx, `UPDATE presences SET expire = $1 WHERE id = ANY ($2)`, now, updates)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		time.Sleep(1 * time.Minute)
	}
}
