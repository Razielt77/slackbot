package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	IN_CHANNEL = "in_channel"
)



func SendSimpleText (url, message string) error {

	msg := slack.Msg{}
	msg.ResponseType = IN_CHANNEL
	msg.Text = message
	_, err := DoPost(url,msg)
	return err
}

func PrintJson(i interface{}){

	bytes, _ := json.Marshal(i)
	fmt.Println(string(bytes))

}

func DoPost (url string, v interface{})([]byte, error){


	jsn, err := json.Marshal(v)

	if err != nil{
		fmt.Println(err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsn))

	if err != nil{
		fmt.Println(err)
		return nil, err
	}

	token := "Bearer " + access_token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	client := &http.Client{}


	resp, err := client.Do(req)

	if err != nil{
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil{
		fmt.Println(err)
		return nil, err
	}

	return body,err

}


func ParseSlashCommand(w http.ResponseWriter, r *http.Request) *slack.SlashCommand {

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return nil
	}

	cmd , err := slack.SlashCommandParse(r)

	if err != nil {
		fmt.Println("Cannot parse %s\n", r.Body)
		http.Error(w, "Cannot Parse", 400)
		return nil
	}

	return &cmd
}


func ComposeLogin(txt string)  *slack.Msg{


	msg := slack.Msg{}
	msg.ResponseType = IN_CHANNEL
	//msg.Text = "*hmm...seems like you haven't logged in recently*"
	msg.Text = txt
	msg.Attachments = []slack.Attachment{*ComposeLoginAttacment()}

	return &msg
}

func ComposeLoginAttacment(msg ...string) *slack.Attachment{
	att := &slack.Attachment{
		Title:"Add your Codefresh's Account Token",
		TitleLink:"https://g.codefresh.io/account-admin/account-conf/tokens#autogen=codefresh-slack-bot",
		Color:"#11b5a4",
		CallbackID: ENTER_TOKEN,
		Text: "Go to your Codefresh's Accounts Settings->Tokens to create your token."}
	if len(msg) == 1{
		att.Title = msg[0]
	}
	att.Actions = []slack.AttachmentAction{{Name: "add-token", Text: "Enter Token", Type: "button",Style:"primary" ,Value: "start"}}
	return att
}

func NormalizeCommit(commit,commit_url string)string{

	commit = strings.Replace(commit,"\n"," ",-1)
	commit = strings.Replace(commit,"\r"," ",-1)
	commit = "<" + commit_url + "|" + commit + ">"
	return commit
}