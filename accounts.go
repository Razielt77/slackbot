package main

import (
	"encoding/json"
	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
	"net/http"
)

func AccountChangeCommand (s *mgo.Session) func(w http.ResponseWriter, r *http.Request){
	return func (w http.ResponseWriter, r *http.Request){


		cmd := ParseSlashCommand(w,r)
		if cmd == nil {
			return
		}

		session := s.Copy()
		defer session.Close()


		usr, _ := GetUser(session,cmd.TeamID,cmd.UserID)


		msg := slack.Msg{}
		msg.ResponseType = "ephemeral"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)

		if usr == nil {

			go DoPost(cmd.ResponseURL, ComposeLogin())
			//composeLogin(&msg)
			return
		}

		//


		//go SendPipelinesListMsg(usr,cmd)

		return

	}
}


