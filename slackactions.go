package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type slackAction struct {
	Name 			string `json:"name"`
	Type			string `json:"type"`
	Value 			string `json:value`
}

type slackActionMsg struct {
	Type 			string `json:"type"`
	CallbackId		string `json:"callback_id"`
	User 			User `json:user`
	Actions			[] slackAction `json:actions`
}

func (r *slackActionMsg) ExtractAction(req *http.Request, log bool) bool {


	err := req.ParseForm()

	if err != nil {
		return false
	}

	payload := req.Form.Get("payload")


	err = json.Unmarshal([]byte(payload), r)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}

	fmt.Printf("User id: %v\nUser name: %v\nActions len: %v\nAction name: %v\nAction type: %v\nAction value: %v\nCallbackid: %v\n", r.User.ID, r.User.Name, len(r.Actions),r.Actions[0].Name,r.Actions[0].Type,r.Actions[0].Value,r.CallbackId)


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

func (r *slackActionMsg) AskToken () bool {

	fmt.Printf("Executing add-token action\n")

	return true
}
