package main

import (
	"context"
	"fmt"
	"guigou/bot-is-gud/api/rpc"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/patrickmn/go-cache"
)

type Presence struct {
	db     *pgx.Conn
	ctx    context.Context
	cache  *cache.Cache
	client rpc.Presence
}

func New(ctx context.Context, db *pgx.Conn) {
	client := rpc.NewPresenceProtobufClient("http://localhost:8080", &http.Client{})
	c := cache.New(2*time.Minute, 2*time.Minute)
	c.OnEvicted(TraceInactive)
	p := &Presence{db: db, ctx: ctx, cache: c, client: client}
	go p.record()
}

func TraceInactive(uid string, id interface{}) {
	fmt.Printf("removing %s active presence\n", uid)
}

func (p *Presence) record() {
	after := time.Now().UTC().Add(-2 * time.Minute)
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
			_, err := p.db.Exec(context.Background(), `UPDATE presences SET expire = $1 WHERE id = ANY ($2)`, now, updates)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		time.Sleep(1 * time.Minute)
	}
}
