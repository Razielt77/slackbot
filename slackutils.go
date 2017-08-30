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


type Action struct {
	Name 	string `json:"name"`
	Text 	string `json:"text"`
	Type 	string `json:"type"`
	Value 	string `json:"value"`
	Style 	string `json:"style"`
}

type Attachment struct {
	Title 				string `json:"title"`
	Callback_id 		string `json:"callback_id"`
	Attachment_type 	string `json:"attachment_type"`
	Actions 			[] Action `json:"actions"`
}




type slackRsp struct {
	Text		string `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

/*

{
	"Text": "*hmm...seems like you haven't logged in recently*",
	"Attachments": [
			{
			"Title": "How would you like to login?",
			"Callback_id": "comic_1234_xyz",
			"Attachment_type": "default",
			"Actions": [
			{
				"Name": "recommend",
				"Text": "Github",
				"Type": "button",
				"Value": "github"
			},
			{
				"Name": "recommend",
				"Text": "Bitbucket",
				"Type": "button",
				"Value": "bitbucket"
			},
			{
				"Name": "recommend",
				"Text": "Gitlab",
				"Type": "button",
				"Value": "gitlab"
			}
		]
		}
	]
}

*/
