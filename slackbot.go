package main

import (
	"fmt"
	"net/http"
	"encoding/json"
)



func handler(w http.ResponseWriter, r *http.Request) {

	var cmd slackCmd

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}


	err := cmd.ExtractCmd(r, true)

	if err != true {
		fmt.Println("Cannot parse %s", r.Body)
		http.Error(w, "Cannot Parse", 400)
		return
	}


	var rsp slackRsp

	rsp.Text = "*hmm...seems like you haven't logged in recently*"
	att := Attachment{Title:"How would you like to login", Callback_id: "login", Attachment_type: "default"}
	att.Actions = []Action{{Name: "github", Text: "Github", Type: "button" ,Value: "github"},
						   {Name: "bitbucket", Text: "Bitbucket", Type: "button" ,Value: "bitbucket"},
						   {Name: "gitlab", Text: "Gitlab", Type: "button" ,Value: "gitlab"}}
	rsp.Attachments = []Attachment{att}


	/*slackRsp = {
		"*hmm...seems like you haven/'t logged in recently*",
		[1]Attachment{"How would you like to login","comic_1234_xyz","default",
			[3]Action{"recommend", "Github","button","github"}
		}

	}*/

	/*if cmd.Text != ""{
		rsp.Text = "All fine. Received text: " + cmd.Text
	}else{
		rsp.Text = "All fine"
	}*/


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rsp)

	//fmt.Fprintf(w, "Hi there, I love %s! and %s", rsp, js)
	//fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)

	/*var rsp slackRsp
	rsp.text = "All fine"
	rsp.id = "dd"
	js, err := json.Marshal(rsp)

	if err != nil {
		fmt.Println("Error")
		return
	}
	fmt.Println("js - %s rsp - %s", js, rsp)*/


	http.ListenAndServe(":8080", nil)
}
