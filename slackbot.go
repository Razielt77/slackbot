package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"os"
)


var access_token string = ""

func handler(w http.ResponseWriter, r *http.Request) {

	var cmd slackCmd
	var rsp slackRsp

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}


    //extracting the command
	err := cmd.ExtractCmd(r, true)

	if err != true {
		fmt.Println("Cannot parse %s", r.Body)
		http.Error(w, "Cannot Parse", 400)
		return
	}

	//if command require login than check if user logged in (have a context) if not it asks him/her to login
	if cmd.LoginRequired() {
		usr, ok := users[cmd.User_id]

		if !ok {
			rsp.composeLogin()
			//users[cmd.User_id] = User{Name:cmd.User_name}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(rsp)
			return
		}
		fmt.Println("User %s run command %s", usr.Name, cmd.Text)

	}

	var clicmd Cfcmd
	if clicmd.ConstructCmd(cmd.Text){

		err, ok := clicmd.RunCmd(&rsp)
		if !ok{
			rsp.Text = "Error executing command err: " + err.Error()
		}
	}else{
		rsp.Text = "Bad command " + cmd.Text
	}



	rsp.ResponseType = "in_channel"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rsp)

}
	var slackApi *slack.Client


func main() {
	//retrieving the slack web api token from the environment variable
	access_token = os.Getenv("TOKEN")




	if access_token == "" || access_token == "not_set"{
		fmt.Printf("WARNING: no access token set value is:%s\n",access_token)
	} else {
		fmt.Printf("Token set is:%s\n",access_token)

	}
	slackApi = slack.New(access_token)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handler)
	router.HandleFunc("/action", handleAction)
	router.HandleFunc("/action", handleAction)
	log.Fatal(http.ListenAndServe(":8080", router))
	//http.HandleFunc("/", handler)
	//http.ListenAndServe(":8080", nil)
}


func handleAction(w http.ResponseWriter, r *http.Request) {


	var action slackActionMsg
	var rsp slackRsp

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	//extracting the command
	err := action.ExecuteAction(r, true)

	if err != true {
		fmt.Printf("Cannot execute %s", r.Body)
		http.Error(w, "Cannot Parse", 400)
		return
	}


	rsp.ResponseType = "in_channel"
	rsp.Text = "In Action baby"


	w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(rsp)

}