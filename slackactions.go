package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type slackAction struct {
	Name 			string `json:"name"`
	Type		string `json:"type"`
	Value 			User `json:value`
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
	}

	fmt.Printf("User id: %v\nUser name: %v\nActions len: %v\nAction: %v\nCallbackid: %v\n", r.User.ID, r.User.Name, len(r.Actions),r.Actions[0],r.CallbackId)

	//fmt.Printf("Command received\n %+v\n", req.Form)




	//fmt.Printf("Payload Type: %T Value: %s\n", payload, payload)
	//fmt.Printf("Type: %T Value: %s\n", r.Type, r.Type)
	//fmt.Printf("Type: %T Value: %s\n", r.CallbackId, r.CallbackId)

	//fmt.Printf("\n\nForm Type: %T Value: %s\n", req.Form, req.Form)


	/*if log != false {

		fmt.Printf("Command received\n %+v\n", r)
	}*/

	return true
}
