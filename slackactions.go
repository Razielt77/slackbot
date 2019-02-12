package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
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

func (r *slackActionMsg) ExecuteAction(req *http.Request, w http.ResponseWriter, log bool) bool {


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
			w.WriteHeader(200)
			SetToken(&intcallback)

		}

	case slack.InteractionTypeInteractionMessage:
			AskToken(&intcallback)

			rsp := slackRsp{ResponseType:"ephemeral",Text:"Prompting token dialog..."}


			//msg := slack.Msg{ResponseType:"ephemeral",Text:"Prompting token dialog..."}
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

func SetToken (callback *slack.InteractionCallback) bool {


	text := ":white_check_mark: *Token submitted!*"
	att := slack.Attachment{
		Color:"#11b5a4",
		Text: "Learn more on Codefresh's slack commands at www.codefresh.io"}

	//msg := slack.Msg{ResponseType:"ephemeral",Text:text,Attachments:[]slack.Attachment{att}}


	fmt.Printf("Using ts to: %s\n", ts )
	//channelID, timestamp, _, err:= slackApi.UpdateMessage(callback.Channel.ID,ts,slack.MsgOptionText(text, false),slack.MsgOptionTS(callback.ActionTs),slack.MsgOptionAttachments(att),slack.MsgOptionUpdate(ts))

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