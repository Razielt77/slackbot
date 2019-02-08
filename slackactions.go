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

func (r *slackActionMsg) ExtractAction(req *http.Request, log bool) bool {


	err := req.ParseForm()

	if err != nil {
		return false
	}

	payload := req.Form.Get("payload")
	fmt.Printf("received payload %s\n", payload)

	err = json.Unmarshal([]byte(payload), r)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}


	//fmt.Printf("User id: %v\nUser name: %v\nActions len: %v\nAction name: %v\nAction type: %v\nAction value: %v\nCallbackid: %v\nTrigger ID: %v\n", r.User.ID, r.User.Name, len(r.Actions),r.Actions[0].Name,r.Actions[0].Type,r.Actions[0].Value,r.CallbackId,r.TriggerID)


	return true
}

func (r *slackActionMsg) ExecuteAction () bool {

	switch r.Actions[0].Name {
	case "add-token":
		if r.AskToken() != true {
			fmt.Println("error asking for token")
			return false
		}
		return true
	default:
		fmt.Printf("Unidentified action %v \n", r.Actions[0].Name)
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

func (r *slackActionMsg) DialogSubmission () bool {

	switch r.CallbackId {
	case "enter_token":
		if r.SetToken() != true {
			fmt.Println("error asking for token")
			return false
		}
		return true
	default:
		fmt.Printf("Unidentified dialog submission %v \n", r.Actions[0].Name)
	}
	return true
}



func (r *slackActionMsg) AskToken () bool {


	fmt.Printf("Executing add-token action\n")


	/*type Dialog struct {
		TriggerID      string          `json:"trigger_id"`      // Required
		CallbackID     string          `json:"callback_id"`     // Required
		State          string          `json:"state,omitempty"` // Optional
		Title          string          `json:"title"`
		SubmitLabel    string          `json:"submit_label,omitempty"`
		NotifyOnCancel bool            `json:"notify_on_cancel"`
		Elements       []DialogElement `json:"elements"`
	}*/

	textElement := &slack.TextInputElement{}
	textElement.Type = "text"
	textElement.Name = "cftoken"
	textElement.Placeholder = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	textElement.Label = "Codefresh Token"

	var dlg slack.Dialog
	dlg.TriggerID = r.TriggerID
	dlg.CallbackID = r.CallbackId
	dlg.Title = "Your Codefresh Token"
	dlg.Elements = []slack.DialogElement{textElement}

	slackApi.OpenDialog(r.TriggerID,dlg)


	return true
}


//func postJSON (url string, )