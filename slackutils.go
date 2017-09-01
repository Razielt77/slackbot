package main

import (
	"net/http"
	"fmt"
	"encoding/json"
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
	Text 				string	`json:"text"`
}




type slackRsp struct {
	Text		string `json:"text"`
	Attachments []Attachment `json:"attachments"`
}


func (r *slackRsp) composeLogin() {


	r.Text = "*hmm...seems like you haven't logged in recently*"
	att := Attachment{Title:"How would you like to login", Callback_id: "login", Attachment_type: "default"}
	att.Actions = []Action{{Name: "login", Text: "Github", Type: "button" ,Value: "github"},
						   {Name: "login", Text: "Bitbucket", Type: "button" ,Value: "bitbucket"},
						   {Name: "login", Text: "Gitlab", Type: "button" ,Value: "gitlab"}}
	r.Attachments = []Attachment{att}

	return
}


func (r *slackRsp) composeLoginScs() {


	att := Attachment{Text: ":white_check_mark: This message  contains a URL <http://foo.com/>", Callback_id: "login", Attachment_type: "default"}
	r.Attachments = []Attachment{att}
	return
}


type actionMsg struct{
	Actions []struct {
						Name string		`json:"name"`
						Value string 	`json:"value"`
						Type string 	`json:"type"`
						} `json:"actions"`
	Callback_id	string	`json:"callback_id"`
	Team		struct{
						Id 		string	`json:"id"`
						Domain 	string	`json:"domain"`
						}	`json:"team"`
	Channel		struct{
		Id 		string	`json:"id"`
		Name 	string	`json:"name"`
	}	`json:"channel"`

	User		struct{
		Id 		string	`json:"id"`
		Name 	string	`json:"name"`
	}	`json:"user"`

	Action_ts			string	`json:"action_ts"`
	Message_ts			string	`json:"message_ts"`
	Attachment_id		string	`json:"attachment_id"`
	Token				string	`json:"token"`
	Original_message	string	`json:"original_message"`
	Response_url		string	`json:"response_url"`

}

func (a *actionMsg) extractMsg (req *http.Request, log bool) bool {
	err := req.ParseForm()

	if err != nil {
		return false
	}

	j := req.Form.Get("payload")

	if log != false {
		fmt.Printf("Command received\n %s\n", j)
	}

	err = json.Unmarshal([]byte(j), a)

	if err != nil {
		fmt.Printf("Cannot unmarshal \n")
		return false
	}

	if log != false {
		fmt.Printf("Marshal succeded for user\n%s\n", a.User.Name)
	}


	return true
}

