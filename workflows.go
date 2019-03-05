package main

import (
	"fmt"
	"github.com/Razielt77/cf-webapi-go"
	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
	"strconv"
	"time"
)

const (
	WORKFLOW_SUCCESS	= "success"
	WORKFLOW_FAIL		=	"error"
	BUILD_URL 			= 	"https://g.codefresh.io/build/"
)

func SendPipelinesWorkflow(s *mgo.Session, intcallback slack.InteractionCallback){


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
	cf_workflows, _ := cfclient.WorkflowList(webapi.OptionLimit("5"),webapi.OptionPipelineID(pipeline_id))

	workflowsMsg := slack.Msg{}

	workflowsMsg.ResponseType = IN_CHANNEL

	workflowsMsg.Text = "*No builds found*"

	if len(cf_workflows) > 0 && err == nil{

		workflowsMsg.Text = "*Showing last " + strconv.Itoa(len(cf_workflows)) + " builds*"
		workflowsMsg.Attachments = ComposeWorkflowAtt(cf_workflows)
	}

	_, err = DoPost(intcallback.ResponseURL,workflowsMsg)

	if err != nil {
		fmt.Printf("Cannot send message\n")
	}
}


func ComposeWorkflowAtt(p_arr []webapi.Workflow) []slack.Attachment {
	var attarr []slack.Attachment

	for _, workflow := range p_arr {
		att := slack.Attachment{ThumbURL: workflow.Avatar}

		field := slack.AttachmentField{
			Title: "Commit",
			Value: NormalizeCommit(workflow.CommitMsg,workflow.CommitUrl),
			Short: false}

		att.Fields = append(att.Fields,field)

		start,duration := ExtractStartAndDuration(workflow.CreatedTS,workflow.FinishedTS)
		field = slack.AttachmentField{
			Title: "Start Time",
			Value: "<!date^" + start + "^{date} at {time}|Not Set>",
			Short:	true}

		att.Fields = append(att.Fields,field)

		var status string
		switch workflow.Status{
		case WORKFLOW_SUCCESS:
			att.Color = "#11b5a4"
			status = ":white_check_mark: Success"
		case WORKFLOW_FAIL:
			att.Color = "#e83f43"
			status = ":heavy_exclamation_mark: Fail"
		default:
			att.Color = "#ccc"
			status = workflow.Status
		}

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

		action := slack.AttachmentAction{
			Text: "Build's logs :spiral_note_pad:",
			URL: BUILD_URL+workflow.ID,
			Type: "button"}

		att.Actions = append(att.Actions,action)

		attarr = append(attarr,att)
	}
	return attarr
}



func ExtractStartAndDuration(start,finish string) (string, string){

	t_start, err := time.Parse(time.RFC3339, start)
	if err != nil {
		fmt.Println(err)
	}

	t_finish, err := time.Parse(time.RFC3339, finish)
	if err != nil {
		fmt.Println(err)
	}
	duration_t := t_finish.Sub(t_start)
	duration := strconv.Itoa(int (duration_t.Minutes())) + " minutes."

	start_unix := strconv.FormatInt(t_finish.Unix(),10)

	return start_unix, duration

}

/*type Workflow struct {
	Status string `json:"status"`
	CreatedTS 	string `json:"created"`
	FinishedTS  string `json:"finished"`
	Committer 	string `json:"userName"`
	CommitMsg 	string `json:"commitMessage"`
	CommitUrl 	string `json:"commitURL"`
	Avatar 		string `json:"avatar"`
}*/


