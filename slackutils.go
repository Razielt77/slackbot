package main

import (
	"net/http"
	"fmt"
)

type slackCmd struct {
	Token    		string      `json:"token"`
	Team_id  		string      `json:"team_id"`
	Team_domain   	string      `json:"team_domain"`
	Enterprise_id	string 		`json:"enterprise_id"`
	Enterprise_name	string		`json:"enterprise_name"`
	Channel_id		string 		`json:"channel_id"`
	Channel_name	string	    `json:"channel_name"`
	User_id			string 		`json:"user_id"`
	User_name		string 		`json:"user_name"`
	Command			string 	  	`json:"command"`
	Text			string		`json:"text"`
	Response_url	string		`json:"response_url"`

}



func (r *slackCmd) ExtractCmd(req *http.Request, log bool) bool {


	err := req.ParseForm()

	if err != nil {
		return false
	}

	r.Token = req.Form.Get("token")
	r.Team_id = req.Form.Get("team_id")
	r.Team_domain = req.Form.Get("team_domain")
	r.Enterprise_id = req.Form.Get("enterprise_id")
	r.Enterprise_name = req.Form.Get("enterprise_name")
	r.Channel_id = req.Form.Get("channel_id")
	r.Channel_name = req.Form.Get("channel_name")
	r.User_id = req.Form.Get("user_id")
	r.User_name = req.Form.Get("user_name")
	r.Command = req.Form.Get("command")
	r.Text = req.Form.Get("text")
	r.Response_url = req.Form.Get("Response_url")

	if log != false {

		fmt.Printf("Command received\n %+v\n", r)
	}

	return true
}



type slackRsp struct {
	Text	string `json:"text"`
}
