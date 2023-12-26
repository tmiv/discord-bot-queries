package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	dgo "github.com/bwmarrin/discordgo"
)

var (
	sess *dgo.Session = nil
)

func authBot() error {
	var err error
	sess, err = dgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	sess.AddHandler(func(s *dgo.Session, r *dgo.Ready) {
		log.Printf("Bot is ready\n")
	})
	sess.Open()
	return err
}

func closeBot() error {
	if sess != nil {
		s := sess
		sess = nil
		return s.Close()
	}
	return nil
}

func verifyMembership(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("Bad Method on verifyMembership %s\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := r.URL.Query().Get("user")
	if len(user) < 2 {
		log.Printf("Bad user on verifyMembership %s\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	channel := r.URL.Query().Get("channel")
	if len(channel) < 2 {
		log.Printf("Bad channel on verifyMembership %s\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	perm, err := sess.UserChannelPermissions(user, channel)
	if err != nil {
		log.Printf("Discord Bot Perm Req Failure %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp := VerifyMembershipResponse{
		UserId:    user,
		ChannelId: channel,
		Member:    (perm & dgo.PermissionViewChannel) != 0,
	}
	encoded, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Marshal erroor %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(encoded)
}
