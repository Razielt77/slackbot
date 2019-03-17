package main

import (
	"fmt"
	"github.com/Razielt77/cf-webapi-go"
	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
	"net/http"
)

const (
	FIRST_TIME_USER = "It seems like your team never added a Codefresh's token. Please add your account's token."
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

func SendTokensList(team *Team, cmd *slack.SlashCommand){


	accountsMsg := slack.Msg{}

	var err error = nil

	PrintJson(team)
	accountsMsg.Text = 	"Currently these are the tokens for your team: " + team.Team


	accountsMsg.Attachments = ComposeTokensAtt(team)
	accountsMsg.Attachments = append(accountsMsg.Attachments,*ComposeAddTokenAtt())
	accountsMsg.ResponseType = IN_CHANNEL

	_, err = DoPost(cmd.ResponseURL,accountsMsg)

	if err != nil {
		fmt.Printf("Cannot send message\n")
	}
}

func ComposeAddTokenAtt() *slack.Attachment {
	att := &slack.Attachment{CallbackID: ENTER_TOKEN}

	att.Actions = []slack.AttachmentAction{{Name: "add-token", Text: "Add Token", Type: "button",Style:"primary" ,Value: "start"}}
	return att
}

func ComposeTokensAtt(team *Team) []slack.Attachment {
	var attarr []slack.Attachment = nil
	for _, account := range team.CFAccounts{
		token := account.Token[len(account.Token)-4:]
		att := slack.Attachment{
			Title: "Account Name: " + account.Name,
			Color:"#11b5a4",
			Text:"Token ends with: " + token + "\nSubmitted by: "+ account.UserName}
		attarr = append(attarr,att)
	}
	return attarr
}

func GetAccountInfo(user *webapi.UserInfo) *webapi.AccountInfo{

	for i, _ := range user.Accounts{
		if user.Accounts[i].Name == user.ActiveAccount {
			return  &user.Accounts[i]
		}
	}
	return nil
}


func AddTokenToTeam(token string, team *Team) error{
	cf_user, err := webapi.New(token).UserInfo()

	if err != nil {
		return err
	}

	err = team.AddToken(GetAccountInfo(cf_user))
	if err != nil {
		return err
	}
	return nil
}

func UpdateTeamTokens (s *mgo.Session, callback *slack.InteractionCallback) bool {

	session := s.Copy()
	defer session.Close()

	team, _ := GetTeam(session,callback.Team.ID)
	token := callback.Submission["cftoken"]

	if team == nil{

		team = &Team{Team:callback.Team.Name,TeamID:callback.Team.ID}

		err := AddTokenToTeam(token,team)
		if err != nil {
			SendSimpleText(callback.ResponseURL,":heavy_exclamation_mark: *Invalid token*: "+ err.Error())
			return false
		}
		fmt.Println("Printing Team before adding...")
		PrintJson(team)
		AddTeam(session,team)


	}else{

		err := AddTokenToTeam(token,team)
		if err != nil {
			SendSimpleText(callback.ResponseURL,":heavy_exclamation_mark: *Invalid token*: "+ err.Error())
			return false
		}

		UpdateTeam(s,team)
	}


	var accountsList string
	for _, account := range team.CFAccounts{
		accountsList = accountsList + "*" + account.Name +"*\n"
	}
	msg := slack.Msg{Text: ":white_check_mark: *Token successfully added!*"}
	att := slack.Attachment{
		Color:"#11b5a4",
		Text: "Enriched URLs are now supported for the following accounts:\n"+ accountsList}

	msg.Attachments = append(msg.Attachments,att)

	DoPost(callback.ResponseURL,msg)


	return true
}



