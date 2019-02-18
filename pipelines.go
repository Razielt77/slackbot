package main

import (
	"encoding/json"
	"github.com/Razielt77/cf-webapi-go"
	"github.com/nlopes/slack"
	"strconv"
	"time"
)

func ComposePipelinesAtt(p_arr []webapi.Pipeline) []slack.Attachment {
	var attarr []slack.Attachment
	for _, pipeline := range p_arr {
		p_att := slack.Attachment{
			Title:pipeline.Name,
			TitleLink:string(`https://g.codefresh.io/pipelines/edit/summary?id=` + pipeline.ID),
			Footer:"Last Executed",
			Color:"#11b5a4"}

		if pipeline.LastWorkflow.CreatedTS != webapi.NO_LAST_WORKFLOW {
			t_start, err := time.Parse(time.RFC3339, pipeline.LastWorkflow.CreatedTS)
			if err != nil {
				panic(err)
			}

			t_finish, err := time.Parse(time.RFC3339, pipeline.LastWorkflow.FinishedTS)
			if err != nil {
				panic(err)
			}
			duration_t := t_finish.Sub(t_start)
			duration := strconv.Itoa(int (duration_t.Minutes())) + " minutes."
			p_att.Ts = json.Number(t_finish.Unix())
			switch pipeline.LastWorkflow.Status{
			case "success":
				p_att.FooterIcon = `https://raw.githubusercontent.com/Razielt77/slackbot/master/img/passed.png`
			case "error":
				p_att.FooterIcon = `https://raw.githubusercontent.com/Razielt77/slackbot/master/img/failed.png`
				p_att.Color ="#e83f43"
			default:
				p_att.Color="#ccc"
			}


			p_att.Fields = append(p_att.Fields,
				slack.AttachmentField{Title:"Last Status", Value:pipeline.LastWorkflow.Status, Short:true},
				slack.AttachmentField{Title:"Duration", Value: duration , Short:true},
				slack.AttachmentField{Title:"Last Commit", Value:pipeline.LastWorkflow.CommitMsg, Short:false})
			p_att.AuthorIcon = pipeline.LastWorkflow.Avatar
			p_att.AuthorName= pipeline.LastWorkflow.Committer
		}
		//status := slack.AttachmentField{}
		attarr = append(attarr, p_att)
	}
	return attarr
}