package main

import (
	"errors"
	"fmt"
	"github.com/Razielt77/cf-webapi-go"
	"github.com/nlopes/slack"
	"regexp"
	"strconv"
	"strings"
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

func SendPipelinesListMsg(usr *User, cmd *slack.SlashCommand){


	msg := slack.Msg{}
	msg.ResponseType = "in_channel"
	msg.Text = "Retrieving Pipelines..."
	DoPost(cmd.ResponseURL,msg)

	//Retrieving the pipelines

	pipelinesMsg := slack.Msg{}

	var err error = nil
	token := usr.GetToken()

	if token == ""{
		fmt.Println("No token found!\n")
	}



	cfclient := webapi.New(token)

	var options []webapi.Option

	options = nil



	if cmd.Text != ""{
		fmt.Printf("optiosn string is:%s\n",cmd.Text)
		options, err = ComposeOption(cmd.Text,TagFlag(),LimitFlag())

		if err != nil {
			fmt.Println(err)
			pipelinesMsg.Text = "*Parsing Error:" + err.Error() +"*"
			DoPost(cmd.ResponseURL,pipelinesMsg)
			return
		}
	}




	pipelines, err := cfclient.PipelinesList(options...)



	pipelinesMsg.Text = "*No Pipelines found*"

	if len(pipelines) > 0 && err == nil{
		pipelinesMsg.Text = "*" + strconv.Itoa(len(pipelines)) + " Pipelines found*"
		pipelinesMsg.Attachments = ComposePipelinesAtt(pipelines)
	}else{
		pipelinesMsg.Text = "*No Pipelines found*"
	}


	resp, err := DoPost(cmd.ResponseURL,pipelinesMsg)

	fmt.Printf("resp is: %s\n",resp)

	//json.NewEncoder(w).Encode(msg)
}

type Flag func() (string, webapi.OptionGen)

func TagFlag() Flag {
	return func() (string, webapi.OptionGen){
		return "tag",webapi.OptionTag
	}
}

func LimitFlag() Flag {
	return func() (string, webapi.OptionGen){
		return "limit",webapi.OptionLimit
	}
}




func ComposeOption(command string, flags ...Flag) ([]webapi.Option , error){

	var options []webapi.Option
	var match_arr []string

	//fmt.Printf("Number of flag types are:%v\n",len(flags))
	for _ , flag := range flags {
		key, option := flag()
		str := `(\s+` + key + `|^tag)\s*?=\s*?\w+`
		re, err := regexp.Compile(str)
		if err != nil{
			fmt.Println(err)
			return nil, errors.New("Parsing error")
		}
		match := re.FindString(command)
		if match != "" {
			arr := strings.Split(match,"=")
			if len(arr) !=2 {
				return nil, errors.New("Parsing error")
			}
			value := strings.Split(match,"=")[1]
			value = strings.Trim(value," ")
			options = append(options, option(value))
			match_arr = append(match_arr,match)
		}

	}

	//checking if there are any redundant characters
	cmd := command
	for _, str := range match_arr {
		cmd = strings.Replace(cmd,str,"",1)
	}
	cmd = strings.Trim(cmd," ")

	if cmd != ""{
		return nil, errors.New("Parsing error")
	}
	return options,nil

}