package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
	"net/http"
)

const (
	FIRST_TIME_USER = "It seems like you never added a Codefresh's token. Please add your account's token."
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

			go DoPost(cmd.ResponseURL, ComposeLogin(FIRST_TIME_USER))
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
	accountsMsg.ResponseType = IN_CHANNEL


	_, err = DoPost(cmd.ResponseURL,accountsMsg)

	if err != nil {
		fmt.Printf("Cannot send message\n")
	}
}

func ComposeAccountsAtt(user *User) []slack.Attachment {
	var attarr []slack.Attachment = nil
	for _, account := range user.CFAccounts{
		att := slack.Attachment{
			Title: "Account Name: " + account.Name,
			Color:"#11b5a4"}

		if account.Name != user.ActiveAccount {
			att.Color = "#ccc"
			if account.Token == "" {
				att.CallbackID = ENTER_TOKEN
				att.Text = "Token required for this account"
				att.Actions = []slack.AttachmentAction{{Name:  "add-token",
														Text:  "Add Token",
														Type:  "button",
														Value: account.Name}}
			}else{
				att.CallbackID = SWITCH_ACCOUNT
				att.Text = "Token exists for this account"
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



