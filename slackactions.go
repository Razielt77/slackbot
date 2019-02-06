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

	user := req.Form.Get("user")


	//fmt.Printf("Command received\n %+v\n", req.Form)

	r.Type = req.Form.Get("type")
	r.CallbackId = req.Form.Get("callback_id")

	fmt.Printf("Type: %T Value: %s\n", user, user)
	fmt.Printf("Type: %T Value: %s\n", r.Type, r.Type)
	fmt.Printf("Type: %T Value: %s\n", r.CallbackId, r.CallbackId)


	/*if log != false {

		fmt.Printf("Command received\n %+v\n", r)
	}*/

	return true
}
