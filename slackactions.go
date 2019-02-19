package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"gopkg.in/mgo.v2"
	"net/http"
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

func (r *slackActionMsg) ExecuteAction(s *mgo.Session,req *http.Request, w http.ResponseWriter, log bool) bool {


	err := req.ParseForm()

	if err != nil {
		return false
	}

	payload := req.Form.Get("payload")
	fmt.Printf("received payload %s\n", payload)


	intcallback := slack.InteractionCallback{}
	err = json.Unmarshal([]byte(payload), &intcallback)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}

	switch intcallback.Type {
	case slack.InteractionTypeDialogSubmission:
		switch intcallback.CallbackID{
		case "enter_token":
			fmt.Printf("token recieved (slack) is: %s\n",intcallback.Submission["cftoken"])

			w.WriteHeader(http.StatusOK)

			w.Header().Set("Content-Type", "application/json")
			SetToken(s, &intcallback)

		}

	case slack.InteractionTypeInteractionMessage:
			AskToken(&intcallback)

			rsp := slackevents.MessageActionResponse{ResponseType:"ephemeral",ReplaceOriginal:true,Text:"Prompting token dialog..."}


			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(rsp)
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

	if user == nil{
		user = &User{TeamID:callback.Team.ID,UserID:callback.User.ID,Name:callback.User.Name,Team:callback.Team.Name,CFTokens:[]CodefreshToken{{AccountName:"Codefresh",Token:callback.Submission["cftoken"],Active:true}}}
		AddUser(session,user)
	}

	text := ":white_check_mark: *Token submitted!*"
	att := slack.Attachment{
		Color:"#11b5a4",
		Text: "Currently supported commands:\ncf-pipelines-list: List pipelines\n"}

	//msg := slack.Msg{ResponseType:"ephemeral",Text:text,Attachments:[]slack.Attachment{att}}


	fmt.Printf("Using ts to: %s\n", ts )

	channelID, timestamp, err := slackApi.PostMessage(callback.Channel.ID, slack.MsgOptionText(text, false),slack.MsgOptionAttachments(att))
	if err != nil {
		fmt.Printf("%s\n", err)
		return false
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)

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
	dlg.CallbackID = callback.CallbackID
	dlg.Title = "Your Codefresh Token"
	dlg.Elements = []slack.DialogElement{textElement}

	ts = callback.ActionTs
	fmt.Printf("Setting ts to: %s\n", ts )

	slackApi.OpenDialog(callback.TriggerID,dlg)


	return true
}

//func postJSON (url string, )