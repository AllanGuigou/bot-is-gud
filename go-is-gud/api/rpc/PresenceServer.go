package rpc

import (
	context "context"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type PresenceServer struct {
	logger *zap.SugaredLogger
	dg     *discordgo.Session
	gid    string
}

func SetupPresenceServer(logger *zap.SugaredLogger, dg *discordgo.Session, gid string) {
	if gid == "" {
		logger.Warnw("failed to setup presence server",
			"error", "invalid guild id",
			"gid", gid)
		return
	}
	g, err := dg.Guild(gid)
	if err != nil || g == nil {
		logger.Warnw("failed to setup presence server",
			"error", err,
			"gid", gid)
		return
	}

	s := &PresenceServer{logger: logger, dg: dg, gid: gid}
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
	g, err := s.State.Guild(gid)
	if err != nil {
		return nil, err
	}

	users := make([]string, 0, len(g.VoiceStates))
	for _, u := range g.VoiceStates {
		if u.ChannelID == g.AfkChannelID {
			continue
		}
		users = append(users, u.UserID)
	}

	return users, nil
}
