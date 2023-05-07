package rpc

import (
	context "context"
	fmt "fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type PresenceServer struct {
	dg  *discordgo.Session
	gid string
}

func SetupPresenceServer(dg *discordgo.Session, gid string) {
	if gid == "" {
		fmt.Println("failed to setup presence server: invalid guild id")
		return
	}
	g, err := dg.Guild(gid)
	if err != nil || g == nil {
		fmt.Printf("failed to setup presence server guild: %v error: %s\n", g, err)
		return
	}

	s := &PresenceServer{dg: dg, gid: gid}
	handler := NewPresenceServer(s)
	go http.ListenAndServe(":8080", handler)
}

func (ps *PresenceServer) WhoseOn(ctx context.Context, req *WhoseOnReq) (*WhoseOnResp, error) {
	users, err := whoseOn(ps.dg, ps.gid, req.VoiceChannel)
	if err != nil {
		return nil, err
	}

	return &WhoseOnResp{Users: users}, nil
}

func whoseOn(s *discordgo.Session, gid, vc string) ([]string, error) {
	fmt.Println("finding whose on...")
	g, err := s.State.Guild(gid)
	if err != nil {
		return nil, err
	}

	users := make([]string, 0, len(g.VoiceStates))
	fmt.Println(len(g.VoiceStates))
	for _, u := range g.VoiceStates {
		if u.ChannelID == g.AfkChannelID {
			continue
		}
		users = append(users, u.UserID)
	}

	return users, nil
}
