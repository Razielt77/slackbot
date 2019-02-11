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

type ActionType struct{
	Type	slack.InteractionType `json:"type"`
}

type slackActionMsg struct {
	Type 			string `json:"type"`
	CallbackId		string `json:"callback_id"`
	User 			User `json:"user"`
	Submission		DialogSubmission `json:"submission"`
	Actions			[] slackAction `json:"actions"`
	TriggerID		string `json:"trigger_id"`
}

func (r *slackActionMsg) ExecuteAction(req *http.Request, log bool) bool {


	err := req.ParseForm()

	if err != nil {
		return false
	}

	payload := req.Form.Get("payload")
	fmt.Printf("received payload %s\n", payload)

	var tp ActionType

	err = json.Unmarshal([]byte(payload), &tp)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}

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
		}

	case slack.InteractionTypeInteractionMessage:
			AskToken(&intcallback)
			}

	err = json.Unmarshal([]byte(payload), r)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}


	return true
}

type TokenSubbmission struct {
	CfToken 	string 	`json:cftoken`
}

func (r *slackActionMsg) SetToken () bool {

	bt, err := json.Marshal(r.Submission)
	if err != nil {
		fmt.Println("error:", err)
	}else{
		fmt.Printf("Submission string is %s\n", string(bt))
	}

	var submission TokenSubbmission
	err = json.Unmarshal([]byte(bt), &submission)
	if err != nil {
		fmt.Println("error:", err)
	}else{
		fmt.Printf("token is %s\n", submission.CfToken)
	}


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

	slackApi.OpenDialog(callback.TriggerID,dlg)


	return true
}

//func postJSON (url string, )