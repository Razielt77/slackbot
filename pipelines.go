package main

import (
	"encoding/json"
	"fmt"
	"github.com/Razielt77/cf-webapi-go"
	"github.com/nlopes/slack"
	"log"
	"strconv"
	"time"
)

const NOT_AVAILABLE string  = "Not Available"

func ComposePipelinesAtt(p_arr []webapi.Pipeline) []slack.Attachment {
	var attarr []slack.Attachment
	for _, pipeline := range p_arr {
		p_att := slack.Attachment{
			Title:pipeline.Name,
			TitleLink:string(`https://g.codefresh.io/pipelines/edit/summary?id=` + pipeline.ID),
			Color:"#ccc",
			Footer: "Last Executed: Not Available"}

		if pipeline.LastWorkflow.Status != webapi.NO_LAST_WORKFLOW {


			t_start, err := time.Parse(time.RFC3339, pipeline.LastWorkflow.CreatedTS)
			if err != nil {
				fmt.Println(err)
			}

			t_finish, err := time.Parse(time.RFC3339, pipeline.LastWorkflow.FinishedTS)
			if err != nil {
				fmt.Println(err)
			}
			duration_t := t_finish.Sub(t_start)
			duration := strconv.Itoa(int (duration_t.Minutes())) + " minutes."


			//p_att.Ts = json.Number(t_finish.Unix())
			switch pipeline.LastWorkflow.Status{
			case "success":
				p_att.FooterIcon = `https://raw.githubusercontent.com/Razielt77/slackbot/master/img/passed.png`
				p_att.Color="#11b5a4"
			case "error":
				p_att.FooterIcon = `https://raw.githubusercontent.com/Razielt77/slackbot/master/img/failed.png`
				p_att.Color ="#e83f43"
			default:
				p_att.Color="#ccc"
			}

			p_att.Footer = "Last Executed: " + "<!date^" + strconv.FormatInt(t_finish.Unix(),10) + "^{date} at {time}|Not Set>"

			commit := "<" + pipeline.LastWorkflow.CommitUrl + "|" + pipeline.LastWorkflow.CommitMsg + ">"
			p_att.Fields = append(p_att.Fields,
				slack.AttachmentField{Title:"Last Status", Value:pipeline.LastWorkflow.Status, Short:true},
				slack.AttachmentField{Title:"Duration", Value: duration , Short:true},
				slack.AttachmentField{Title:"Last Commit", Value:commit, Short:false})
			p_att.AuthorIcon = pipeline.LastWorkflow.Avatar
			p_att.AuthorName= pipeline.LastWorkflow.Committer
		}else{
			p_att.Fields = append(p_att.Fields,
				slack.AttachmentField{Title:"Last Status", Value:NOT_AVAILABLE, Short:true},
				slack.AttachmentField{Title:"Duration", Value: NOT_AVAILABLE , Short:true},
				slack.AttachmentField{Title:"Last Commit", Value:NOT_AVAILABLE, Short:false})
			p_att.AuthorName= pipeline.LastWorkflow.Committer
		}


		attarr = append(attarr, p_att)
	}
	return attarr
}

func SendPipelinesListMsg(usr *User, response_url string){
	//Retrieving the pipelines

	pipelinesMsg := slack.Msg{}


	token := usr.GetToken()

	if token == ""{
		fmt.Println("No token found!\n")
	}



	cfclient := webapi.New(token)

	pipelines, err := cfclient.PipelinesList()

	log.Printf("Token used:%s\nNumber of pipelines is:%v\n",token,len(pipelines))


	pipelinesMsg.Text = "*No Pipelines found*"

	if len(pipelines) > 0 && err == nil{
		pipelinesMsg.Text = "*" + strconv.Itoa(len(pipelines)) + " Pipelines found*"

		att_arr := ComposePipelinesAtt(pipelines)
		str, err := json.Marshal(att_arr)
		fmt.Printf("att_arr is: %s\n",str)

		if err != nil {
			fmt.Println(err)
			return
		}

		pipelinesMsg.Attachments = att_arr
	}else{
		pipelinesMsg.Text = "*No Pipelines found*"
	}


	resp, err := DoPost(response_url,pipelinesMsg)

	fmt.Printf("resp is: %s\n",resp)

	//json.NewEncoder(w).Encode(msg)
}
