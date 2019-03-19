package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
	"net/http"
)



type slackAction struct {
	Name 			string `json:"name"`
	Type			string `json:"type"`
	Value 			string `json:"value"`
}


type DialogSubmission interface {}



type slackActionMsg struct {
	Type 			string `json:"type"`
	CallbackId		string `json:"callback_id"`
	User 			User `json:"user"`
	Submission		DialogSubmission `json:"submission"`
	Actions			[] slackAction `json:"actions"`
	TriggerID		string `json:"trigger_id"`
}

//var ts string

type slackUnfurlActionResponse struct {
	IsAppUnfurl	bool	`json:"is_app_unfurl"`
}



func HandleAction (s *mgo.Session) func(w http.ResponseWriter, r *http.Request){
	return func (w http.ResponseWriter, r *http.Request){

		w.WriteHeader(http.StatusOK)
		session := s.Copy()
		defer session.Close()

		var action slackActionMsg
		//var rsp slackRsp

		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		action.ExecuteAction(s,r, true)

	}
}

func (r *slackActionMsg) ExecuteAction(s *mgo.Session,req *http.Request, log bool) bool {

	session := s.Copy()
	defer session.Close()

	err := req.ParseForm()

	if err != nil {
		return false
	}

	payload := req.Form.Get("payload")
	fmt.Printf("received payload %s\n", payload)


	intcallback := slack.InteractionCallback{}

	unfurlResponse := slackUnfurlActionResponse{IsAppUnfurl:false}
	err = json.Unmarshal([]byte(payload), &unfurlResponse)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}

	if unfurlResponse.IsAppUnfurl{
		fmt.Println("Response received from unfurl action")
		//TODO currently ignoring error in unmarshaling an unfurl action message - atachment id received as string instead of int
		json.Unmarshal([]byte(payload), &intcallback)
	}else{
		err = json.Unmarshal([]byte(payload), &intcallback)
		if err != nil {
			fmt.Println("error:", err)
			return false
		}
	}


	switch intcallback.Type {
	case slack.InteractionTypeDialogSubmission:
		switch intcallback.CallbackID{
		case ENTER_TOKEN:
			fmt.Printf("token recieved (slack) is: %s\n",intcallback.Submission["cftoken"])
			//w.WriteHeader(http.StatusOK)
			//SetToken(s, &intcallback)
			UpdateTeamTokens(s,&intcallback)
		}

	case slack.InteractionTypeInteractionMessage:

		fmt.Printf("Type is: %s\nCallback ID is: %s\n",intcallback.Type,intcallback.CallbackID)

		switch intcallback.CallbackID {
		case SWITCH_ACCOUNT:
			SwitchAccount(session, &intcallback)
			go slackApi.DeleteMessage(intcallback.Channel.ID, intcallback.MessageTs)
		case ENTER_TOKEN:
			//w.WriteHeader(http.StatusOK)
			AskToken(&intcallback)
			go slackApi.DeleteMessage(intcallback.Channel.ID, intcallback.MessageTs)
		case PIPELINE_ACTION:
			switch intcallback.Actions[0].Name {
			case VIEW_BUILDS:
				go SendPipelinesWorkflow(s,&intcallback)
				//SendSimpleText(intcallback.ResponseURL, "Asking to view builds for "+intcallback.Actions[0].Value)
				}
			}
		}

	err = json.Unmarshal([]byte(payload), r)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}


	return true
}




func SwitchAccount (s *mgo.Session, callback *slack.InteractionCallback) bool {

	//user := User{TeamID:callback.Team.ID,UserID:callback.User.ID,Name:callback.User.Name,Team:callback.Team.Name}
	session := s.Copy()
	defer session.Close()

	user, _ := GetUser(session,callback.Team.ID,callback.User.ID)

	if user == nil{
		SendSimpleText(callback.ResponseURL,"User not exist!")
	}else{
		user.ActiveAccount = callback.Actions[0].Value
		UpdateUser(s,user)
	}

	msg := slack.Msg{Text: ":white_check_mark: *Account switched successfully!*"}
	att := slack.Attachment{
		Color:"#11b5a4",
		Text: "Active account is: *" +
			user.ActiveAccount +
			"*\nCurrently supported commands*:\n" +
			"*/cf-pipelines-list*  Lists pipelines.\n"+
			"*/cf-pipelines-list-active*  Lists pipelines active past week.\n" +
			"*/cf-switch-account* Switch between your Codefresh's accounts.\n",
		ThumbURL: user.Avatar}

	msg.Attachments = append(msg.Attachments,att)

	DoPost(callback.ResponseURL,msg)


	return true
}



func AskToken (callback *slack.InteractionCallback) bool {


	//fmt.Printf("Executing add-token action\n")

	textElement := &slack.TextInputElement{}
	textElement.Type = "text"
	textElement.Name = "cftoken"
	textElement.Placeholder = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	textElement.Label = "Codefresh Token"

	var dlg slack.Dialog
	dlg.TriggerID = callback.TriggerID
	//storing the url in the state
	dlg.State = callback.Actions[0].Value
	dlg.CallbackID = callback.CallbackID
	dlg.Title = "Your Codefresh Token"
	dlg.Elements = []slack.DialogElement{textElement}

	//ts = callback.ActionTs
	//fmt.Printf("Setting ts to: %s\n", ts )

	slackApi.OpenDialog(callback.TriggerID,dlg)


	return true
}

//func postJSON (url string, )