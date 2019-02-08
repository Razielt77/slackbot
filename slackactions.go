package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type cftokenDialogElement struct {
	Label			string `json:"label"`
	Name			string `json:"name"`
	Type 			string `json:"type"`
	Placeholder 	string `json:"placeholder"`
}


type cftokenDialog struct {
	CallbackID		string `json:"callback_id"`
	Title			string `json:"title"`
	SubmitLabel 	string `json:"submit_label"`
	Elements 		[]cftokenDialogElement `json:"elements"`
}

type cftokenDialogMsg struct {
	TriggerID		string `json:"trigger_id"`
	Dialog			cftokenDialog `json:"dialog"`
}

func (r *slackActionMsg) AskToken () bool {


	fmt.Printf("Executing add-token action\n")

	var tknDlg cftokenDialogMsg
	tknDlg.TriggerID = r.TriggerID
	tknDlg.Dialog.CallbackID = r.CallbackId
	tknDlg.Dialog.Title = "Your Codefresh Token"
	tknDlg.Dialog.Elements = []cftokenDialogElement{{Name: "cftoken", Label: "Codefresh Token", Type: "text",Placeholder: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"}}

	bearer := "Bearer " + access_token
	url := "https://slack.com/api/dialog.open"


	bt, err := json.Marshal(tknDlg)
	if err != nil {
		fmt.Println("error:", err)
	}


	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bt))
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	fmt.Printf("Sending Dialog Json: %s\n", string(bt))

	fmt.Printf("Sending to URL: %s, Header: %v\n", req.URL, req.Header)


	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
	}


	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err ==nil {
			bodyString := string(bodyBytes)
			fmt.Printf("Received Response: Status: %s, Body: %s\n", resp.Status, bodyString)
		}else{
			fmt.Println("error:", err)
		}
	}else{
		fmt.Printf("Received Response: Status: %s\n", resp.Status)
	}



	return true
}


//func postJSON (url string, )