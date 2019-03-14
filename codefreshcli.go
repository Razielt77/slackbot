package main

import (
	"github.com/Razielt77/slack"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

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

func (cmd *Cfcmd) RunCmd (rsp *slack.Msg) (err error, ok bool){
	var out []byte
	switch cmd.command {
	case "version":
		out, err = exec.Command("codefresh", "version").Output()
		if err != nil {
			log.Fatal(err)
			return err, false
		}
		re := regexp.MustCompile(`[0-9]*\.[0-9]*\.[0-9]*`)
		cliVer := re.FindString(string(out))
		rsp.Text = "*Slackbot version: 0.0.1*\n*Codefresh CLI version: " + cliVer +"*"
	default:
		rsp.Text = "*" + cmd.command + " is not supported yet. Stay tune.*"

	}

	return nil, true
}

