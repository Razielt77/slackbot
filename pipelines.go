package main

import (
	"encoding/json"
	"github.com/Razielt77/cf-webapi-go"
	"github.com/nlopes/slack"
	"time"
)

func ComposePipelinesAtt(p_arr []webapi.Pipeline) []slack.Attachment {
	var attarr []slack.Attachment
	for _, pipeline := range p_arr {
		p_att := &slack.Attachment{
			Title:pipeline.Name,
			TitleLink:string(`https://g.codefresh.io/pipelines/edit/summary?id=` + pipeline.ID),
			Footer:"Last Executed",
			Color:"#11b5a4"}

		if pipeline.LastWorkflow.CreatedTS != "N/A" {
			t, err := time.Parse(time.RFC3339, pipeline.LastWorkflow.CreatedTS)
			if err != nil {
				panic(err)
			}
			p_att.Ts = json.Number(t.Unix())
			switch pipeline.LastWorkflow.Status{
			case "success":
				p_att.FooterIcon = `https://raw.githubusercontent.com/Razielt77/slackbot/master/img/passed.png`
			case "error":
				p_att.FooterIcon = `https://raw.githubusercontent.com/Razielt77/slackbot/master/img/failed.png`
			}
		}

		//status := slack.AttachmentField{}
		//attarr = append(attarr, p_att)
	}
	return attarr
}
