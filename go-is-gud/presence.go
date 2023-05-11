package main

import (
	"context"
	"fmt"
	"guigou/bot-is-gud/api/rpc"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
)

type Presence struct {
	db     *pgx.Conn
	ctx    context.Context
	client rpc.Presence
}

func New(ctx context.Context, db *pgx.Conn) {
	client := rpc.NewPresenceProtobufClient("http://localhost:8080", &http.Client{})
	p := &Presence{db: db, ctx: ctx, client: client}
	go p.record()
}

func (p *Presence) record() {
	for {
		res, err := p.client.WhoseOn(p.ctx, &rpc.WhoseOnReq{VoiceChannel: ""})
		if err != nil {
			fmt.Println(err)
			return
		}

		rows, err := p.db.Query(context.Background(), `SELECT id, uid, expire FROM presences WHERE active = TRUE`)
		if err != nil {
			fmt.Println(err)
			return
		}

		active := make(map[string]int64)
		for rows.Next() {
			var id int64
			var uid string
			// TODO: how to handle gap in availability where presences may not have been updated for a while
			var expire time.Time
			err = rows.Scan(&id, &uid, &expire)
			if err != nil {
				fmt.Println(err)
				return
			}

			active[uid] = id
		}

		for _, u := range res.Users {
			id, found := active[u]
			delete(active, u)
			if found {
				fmt.Printf("found %s active presence\n", u)
				_, err := p.db.Exec(context.Background(), `UPDATE presences SET expire = $1 WHERE id = $2`, time.Now().UTC(), id)
				if err != nil {
					fmt.Println(err)
					return
				}
			} else {
				fmt.Printf("created %s active presence\n", u)
				_, err := p.db.Exec(context.Background(), `INSERT INTO presences (uid, active, start, expire) VALUES ($1, TRUE, $2, $3)`, u, time.Now().UTC(), time.Now().UTC())
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		ids := make([]int64, 0)
		for uid, id := range active {
			fmt.Printf("found %s inactive presence\n", uid)
			ids = append(ids, id)
		}

		if len(ids) > 0 {
			_, err := p.db.Exec(context.Background(), `UPDATE presences SET active = FALSE WHERE id = ANY ($1)`, ids)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		time.Sleep(1 * time.Minute)
	}
}
