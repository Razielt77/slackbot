package main

import (
	"encoding/json"
	"fmt"
	"github.com/Razielt77/cf-webapi-go"
	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
	"net/http"
)

const (
	ENTER_TOKEN = "enter_token"
	SWITCH_ACCOUNT = "switch_account"
)

type slackAction struct {
	Name 			string `json:"name"`
	Type			string `json:"type"`
	Value 			string `json:"value"`
}


type DialogSubmission interface {}



type slackActionMsg struct {
	Type 			string `json:"type"`
	CallbackId		string `json:"callback_id"`
	User 			User `json:"user"`
	Submission		DialogSubmission `json:"submission"`
	Actions			[] slackAction `json:"actions"`
	TriggerID		string `json:"trigger_id"`
}

var ts string

type slackRsp struct {
	ResponseType string `json:"response_type"`
	Text		string `json:"text"`
	Attachments []slack.Attachment `json:"attachments"`
}

func (r *slackActionMsg) ExecuteAction(s *mgo.Session,req *http.Request, log bool) bool {

	session := s.Copy()
	defer session.Close()

	err := req.ParseForm()

	if err != nil {
		return false
	}

	payload := req.Form.Get("payload")
	//fmt.Printf("received payload %s\n", payload)


	intcallback := slack.InteractionCallback{}


	err = json.Unmarshal([]byte(payload), &intcallback)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}


	switch intcallback.Type {
	case slack.InteractionTypeDialogSubmission:
		switch intcallback.CallbackID{
		case ENTER_TOKEN:
			fmt.Printf("token recieved (slack) is: %s\n",intcallback.Submission["cftoken"])
			//w.WriteHeader(http.StatusOK)
			SetToken(s, &intcallback)
		}

	case slack.InteractionTypeInteractionMessage:

		fmt.Printf("Type is: %s\nCallback ID is: %s\n",intcallback.Type,intcallback.CallbackID)

		switch intcallback.CallbackID {
		case SWITCH_ACCOUNT:
			SwitchAccount(session, &intcallback)
			go slackApi.DeleteMessage(intcallback.Channel.ID, intcallback.MessageTs)
		case ENTER_TOKEN:
			//w.WriteHeader(http.StatusOK)
			AskToken(&intcallback)
			go slackApi.DeleteMessage(intcallback.Channel.ID, intcallback.MessageTs)
		case PIPELINE_ACTION:
			switch intcallback.Actions[0].Name {
			case VIEW_BUILDS:
				go SendPipelinesWorkflow(s,intcallback)
				//SendSimpleText(intcallback.ResponseURL, "Asking to view builds for "+intcallback.Actions[0].Value)
				}
			}
		}

	err = json.Unmarshal([]byte(payload), r)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}


	return true
}

func SetToken (s *mgo.Session, callback *slack.InteractionCallback) bool {

	//user := User{TeamID:callback.Team.ID,UserID:callback.User.ID,Name:callback.User.Name,Team:callback.Team.Name}
	session := s.Copy()
	defer session.Close()

	user, _ := GetUser(session,callback.Team.ID,callback.User.ID)
	token := callback.Submission["cftoken"]

	if user == nil{


		user = &User{TeamID:callback.Team.ID,UserID:callback.User.ID,Name:callback.User.Name,Team:callback.Team.Name}

		//retrieving user's accounts

		cf_user, err := webapi.New(token).UserInfo()

		if err != nil {
			SendSimpleText(callback.ResponseURL,":heavy_exclamation_mark: *Invalid token*: "+ err.Error())
			return false
		}

		user.CFUserName = cf_user.Name
		user.CFAccounts = cf_user.Accounts
		user.ActiveAccount = cf_user.ActiveAccount
		user.Avatar = cf_user.UserData.Image
		err = user.SetToken(token)

		if err != nil{
			fmt.Println(err)
		}

		AddUser(session,user)

	}else{
		//fmt.Printf("Token submitted to exisiting account")

		//checking that tken is valid
		cf_user, err := webapi.New(token).UserInfo()

		if err != nil {
			SendSimpleText(callback.ResponseURL,":heavy_exclamation_mark: *Invalid token*: "+ err.Error())
			return false
		}

		if(cf_user.ActiveAccount != callback.State){
			SendSimpleText(callback.ResponseURL,":heavy_exclamation_mark: *Invalid token*: Token doesn't match with selected account")
			return false
		}

		user.ActiveAccount = callback.State
		user.SetToken(token)
		UpdateUser(s,user)
	}

	msg := slack.Msg{Text: ":white_check_mark: *Token successfully submitted!*"}
	att := slack.Attachment{
		Color:"#11b5a4",
		Text: "Welcome *"+user.CFUserName +
			  "!*\nActive account is: *" +
			   user.ActiveAccount +
			  "*\nCurrently supported commands:\n" +
			  "*/cf-pipelines-list*  Lists pipelines.\n"+
			  "*/cf-pipelines-list-active*  Lists pipelines active past week.\n" +
			  "*/cf-switch-account* Switch between your Codefresh's accounts.\n",
		ThumbURL: user.Avatar}

	msg.Attachments = append(msg.Attachments,att)

	DoPost(callback.ResponseURL,msg)


	return true
}


func SwitchAccount (s *mgo.Session, callback *slack.InteractionCallback) bool {

	//user := User{TeamID:callback.Team.ID,UserID:callback.User.ID,Name:callback.User.Name,Team:callback.Team.Name}
	session := s.Copy()
	defer session.Close()

	user, _ := GetUser(session,callback.Team.ID,callback.User.ID)

	if user == nil{
		SendSimpleText(callback.ResponseURL,"User not exist!")
	}else{
		user.ActiveAccount = callback.Actions[0].Value
		UpdateUser(s,user)
	}

	msg := slack.Msg{Text: ":white_check_mark: *Account switched successfully!*"}
	att := slack.Attachment{
		Color:"#11b5a4",
		Text: "Active account is: *" +
			user.ActiveAccount +
			"*\nCurrently supported commands*:\n" +
			"*/cf-pipelines-list*  Lists pipelines.\n"+
			"*/cf-pipelines-list-active*  Lists pipelines active past week.\n" +
			"*/cf-switch-account* Switch between your Codefresh's accounts.\n",
		ThumbURL: user.Avatar}

	msg.Attachments = append(msg.Attachments,att)

	DoPost(callback.ResponseURL,msg)


	return true
}





func AskToken (callback *slack.InteractionCallback) bool {


	fmt.Printf("Executing add-token action\n")

	textElement := &slack.TextInputElement{}
	textElement.Type = "text"
	textElement.Name = "cftoken"
	textElement.Placeholder = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	textElement.Label = "Codefresh Token"

	var dlg slack.Dialog
	dlg.TriggerID = callback.TriggerID
	//storing the desired account in the state
	dlg.State = callback.Actions[0].Value
	dlg.CallbackID = callback.CallbackID
	dlg.Title = "Your Codefresh Token"
	dlg.Elements = []slack.DialogElement{textElement}

	ts = callback.ActionTs
	fmt.Printf("Setting ts to: %s\n", ts )

	slackApi.OpenDialog(callback.TriggerID,dlg)


	return true
}

//func postJSON (url string, )