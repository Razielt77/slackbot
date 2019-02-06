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

	fmt.Printf("Command received\n %+v\n", req.Form)

	r.Type = req.Form.Get("type")
	r.CallbackId = req.Form.Get("callback_id")


	/*if log != false {

		fmt.Printf("Command received\n %+v\n", r)
	}*/

	return true
}
