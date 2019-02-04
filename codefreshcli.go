package main

import "strings"

type Cfcmd struct {

	command string
	arg []string
}

func (cmd *Cfcmd) ConstructCmd (str string) bool {

	cmd.arg = strings.Split(str, " ")

	if len(cmd.arg) < 1 {
		return false
	}
	cmd.command = cmd.arg[0]

	//if len(cmd.arg) < 1 || cmd.arg[0] != "codefresh" {
	//	return false
	//}
	//cmd.command = cmd.arg[0]
	//cmd.arg = cmd.arg[1:]

	return true
}

func (cmd *Cfcmd) RunCmd (rsp *slackRsp) (err string, ok bool){

	rsp.Text = "Slackbot version: 0.1\n Codefresh CLI version: 2.2"

	return "", true
}

