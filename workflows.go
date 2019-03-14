package main

import (
	"fmt"
	"github.com/Razielt77/cf-webapi-go"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"gopkg.in/mgo.v2"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	WORKFLOW_SUCCESS	= "success"
	WORKFLOW_FAIL		=	"error"
	WORKFLOW_RUNNING		=	"running"
	BUILD_URL 			= 	"https://g.codefresh.io/build/"
)

func SendPipelinesWorkflow(s *mgo.Session, intcallback *slack.InteractionCallback){

	go slackApi.DeleteMessage(intcallback.Channel.ID, intcallback.MessageTs)

	if SendSimpleText(intcallback.ResponseURL,"Retrieving Workflow...") != nil {
		fmt.Printf("Cannot send message\n")
		return
	}

	session := s.Copy()
	defer session.Close()

	usr, _ := GetUser(session,intcallback.Team.ID,intcallback.User.ID)

	if usr == nil {
		SendSimpleText(intcallback.ResponseURL,"Internal error no user found...")
		return
	}


	var err error = nil
	token := usr.GetToken()

	if token == ""{
		SendSimpleText(intcallback.ResponseURL,"Internal error: no token found...")
		return
	}

	cfclient := webapi.New(token)
	pipeline_id := intcallback.Actions[0].Value
	pipline_arr, _ := cfclient.PipelinesList(webapi.OptionID(pipeline_id))

	if len(pipline_arr) != 1{
		SendSimpleText(intcallback.ResponseURL,"Internal error: cant get pipeline...")
		return
	}
	cf_workflows, _ := cfclient.WorkflowList(webapi.OptionLimit("5"),webapi.OptionPipelineID(pipeline_id))

	workflowsMsg := slack.Msg{}

	workflowsMsg.ResponseType = IN_CHANNEL

	workflowsMsg.Text = "*No builds found*"

	if len(cf_workflows) > 0 && err == nil{

		workflowsMsg.Text = "Showing last " +
			strconv.Itoa(len(cf_workflows)) +
			" builds for pipeline: *" + pipline_arr[0].Name +
			"*"
		workflowsMsg.Attachments = ComposeWorkflowAttArr(cf_workflows)
	}

	_, err = DoPost(intcallback.ResponseURL,workflowsMsg)

	if err != nil {
		fmt.Printf("Cannot send message\n")
	}
}


func EnrichSharedLink(s *mgo.Session, team_id string, event *slackevents.LinkSharedEvent){

	session := s.Copy()
	defer session.Close()
	m := make(map[string]slack.Attachment)
	var att *slack.Attachment

	url := event.Links[0].URL
	build := ExtractBuildFromURL(url)

	if build == ""{
		return
	}

	//retrieving user
	usr, _ := GetUser(session,team_id,event.User)


	if usr == nil {
		att = ComposeLoginAttacment("Add Codefresh's token for enriched link messages")
		m[event.Links[0].URL] = *att
		//slackApi.UnfurlMessage(event.Channel,event.MessageTimeStamp.String(),m)
		return
	}

	workflow := GetWorkflowInfo(build,usr)

	if workflow == nil {
		att = ComposeLoginAttacment("Token for this account wan't submitted. Add Token to view enriched link for this account")
		m[event.Links[0].URL] = *att
		//slackApi.UnfurlMessage(event.Channel,event.MessageTimeStamp.String(),m)
		return
	}

	att = ComposeWorkflowAttachment(workflow)

	m[event.Links[0].URL] = *att

	slackApi.UnfurlMessage(event.Channel,event.MessageTimeStamp.String(),m)

}

func GetWorkflowInfo(id string, usr *User) *webapi.Workflow{

	var workflow *webapi.Workflow = nil

	for _, account := range usr.CFAccounts{
		if account.Token != ""{
			client := webapi.New(account.Token)
			workflow, _ = client.GetBuild(id)
			if workflow != nil {return workflow}
		}
	}
	return workflow
}


func ExtractBuildFromURL(url string) (build string){
	workflow_url_reg, _ := regexp.Compile(BUILD_URL+`[0-9A-Za-z]+`)

	if workflow_url_reg.MatchString(url){
		build = strings.TrimLeft(url,BUILD_URL)
	}else{
		build = ""
	}
	return
}



func ComposeWorkflowAttArr(p_arr []webapi.Workflow) []slack.Attachment {
	var attarr []slack.Attachment

	for _, workflow := range p_arr {
		attarr = append(attarr,*ComposeWorkflowAttachment(&workflow))
	}
	return attarr
}


func ComposeWorkflowAttachment(workflow *webapi.Workflow) *slack.Attachment{

	att := &slack.Attachment{ThumbURL: workflow.Avatar}

	att.Title = workflow.Project

	field := slack.AttachmentField{
		Title: "Commit",
		Value: NormalizeCommit(workflow.CommitMsg,workflow.CommitUrl),
		Short: false}

	att.Fields = append(att.Fields,field)

	finish_ts := workflow.FinishedTS




	var status string
	switch workflow.Status{
	case WORKFLOW_SUCCESS:
		att.Color = "#11b5a4"
		status = ":white_check_mark: Success"
	case WORKFLOW_FAIL:
		att.Color = "#e83f43"
		status = ":heavy_exclamation_mark: Fail"
	case WORKFLOW_RUNNING:
		att.Color = "#6AA9DA"
		status = ":gear: Running"
		finish_ts = ""
	default:
		att.Color = "#ccc"
		status = workflow.Status
	}

	start,duration := ExtractStartAndDuration(workflow.CreatedTS,finish_ts)
	field = slack.AttachmentField{
		Title: "Start Time",
		Value: "<!date^" + start + "^{date} at {time}|Not Set>",
		Short:	true}

	att.Fields = append(att.Fields,field)

	field = slack.AttachmentField{
		Title: "Status",
		Value: status,
		Short: true}

	att.Fields = append(att.Fields,field)

	field = slack.AttachmentField{
		Title: "Duration",
		Value: duration,
		Short: true}

	att.Fields = append(att.Fields,field)

	field = slack.AttachmentField{
		Title: "Branch",
		Value: workflow.Branch,
		Short: true}

	att.Fields = append(att.Fields,field)

	field = slack.AttachmentField{
		Title: "SHA",
		Value: workflow.SHA,
		Short: false}

	att.Fields = append(att.Fields,field)

	action := slack.AttachmentAction{
		Text: "Build's logs :spiral_note_pad:",
		URL: BUILD_URL+workflow.ID,
		Type: "button"}

	att.Actions = append(att.Actions,action)
	return att
}



func ExtractStartAndDuration(start,finish string) (string, string){

	t_start, err := time.Parse(time.RFC3339, start)
	if err != nil {
		fmt.Println(err)
	}

	var t_finish time.Time
	if finish != ""{
		t_finish, err = time.Parse(time.RFC3339, finish)
		if err != nil {
			fmt.Println(err)
		}
	}else{
		t_finish = time.Now()
	}

	duration_t := t_finish.Sub(t_start)
	duration := strconv.Itoa(int (duration_t.Minutes())) + " minutes."

	start_unix := strconv.FormatInt(t_finish.Unix(),10)

	return start_unix, duration

}




