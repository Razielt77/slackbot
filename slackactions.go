package main

import (
	"fmt"
	"net/http"
)

type slackActionUser struct {
	Id string `json:"id"`
	Name		string `json:"name"`
}
type slackAction struct {
	Type 			string `json:"type"`
	CallbackId		string `json:"callback_id"`
	User 			slackActionUser `json:user`
}

func (r *slackAction) ExtractAction(req *http.Request, log bool) bool {


	err := req.ParseForm()

	if err != nil {
		return false
	}

	payload := req.Form.Get("payload")


	//fmt.Printf("Command received\n %+v\n", req.Form)

	r.Type = req.Form.Get("type")
	r.CallbackId = req.Form.Get("callback_id")


	fmt.Printf("Payload Type: %T Value: %s\n", payload, payload)
	//fmt.Printf("Type: %T Value: %s\n", r.Type, r.Type)
	//fmt.Printf("Type: %T Value: %s\n", r.CallbackId, r.CallbackId)

	//fmt.Printf("\n\nForm Type: %T Value: %s\n", req.Form, req.Form)


	/*if log != false {

		fmt.Printf("Command received\n %+v\n", r)
	}*/

	return true
}
