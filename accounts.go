package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
	"net/http"
)

func AccountChangeCommand (s *mgo.Session) func(w http.ResponseWriter, r *http.Request){
	return func (w http.ResponseWriter, r *http.Request){


		w.WriteHeader(http.StatusOK)

		cmd := ParseSlashCommand(w,r)
		if cmd == nil {
			return
		}

		session := s.Copy()
		defer session.Close()

		usr, _ := GetUser(session,cmd.TeamID,cmd.UserID)

		if usr == nil {

			go DoPost(cmd.ResponseURL, ComposeLogin())
			return
		}

		go SendAccountsList(usr,cmd)

		return

	}
}


func SendAccountsList(usr *User, cmd *slack.SlashCommand){


	accountsMsg := slack.Msg{}

	var err error = nil

	accountsMsg.Text = 	"Current default account: *" + usr.ActiveAccount + "*\n"+
		 				"Choose one of the below accounts to switch active account"


	accountsMsg.Attachments = ComposeAccountsAtt(usr)


	_, err = DoPost(cmd.ResponseURL,accountsMsg)

	if err != nil {
		fmt.Printf("Cannot send message\n")
	}
}

func ComposeAccountsAtt(user *User) []slack.Attachment {
	var attarr []slack.Attachment = nil
	for _, account := range user.CFAccounts{
		att := slack.Attachment{
			Title:account.Name,
			Color:"#ccc"}

		if account.Name != user.ActiveAccount {
			if account.Token == "" {
				att.CallbackID = ENTER_TOKEN
				att.Actions = []slack.AttachmentAction{{Name:  "add-token",
														Text:  "Add Token",
														Type:  "button",
														Style: "primary" ,
														Value: account.Name}}
			}else{
				att.CallbackID = SWITCH_ACCOUNT
				att.Actions = []slack.AttachmentAction{{Name:  "add-token",
					Text:  "Set Active",
					Type:  "button",
					Style: "primary" ,
					Value: account.Name}}
			}
		}else {
			att.Text = "*Currently Active*"
		}
		attarr = append(attarr,att)
	}
	return attarr
}



