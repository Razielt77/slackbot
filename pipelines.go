package main

/*func ComposePipelinesAtt(p_arr []webapi.Pipeline) []slack.Attachment {
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
		}

		//status := slack.AttachmentField{}
		//attarr = append(attarr, p_att)
	}
	return attarr
}*/
