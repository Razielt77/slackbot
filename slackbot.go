package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"os"
	"time"
)


var access_token string = ""

func handler(w http.ResponseWriter, r *http.Request) {

	var cmd slackCmd
	//var rsp slackRsp

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}


    //extracting the command
	err := cmd.ExtractCmd(r, true)

	if err != true {
		fmt.Println("Cannot parse %s\n", r.Body)
		http.Error(w, "Cannot Parse", 400)
		return
	}

	msg := slack.Msg{}
	//if command require login than check if user logged in (have a context) if not it asks him/her to login
	if cmd.LoginRequired() {
		usr, ok := users[cmd.User_id]

		if !ok {
			composeLogin(&msg)
			//users[cmd.User_id] = User{Name:cmd.User_name}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(msg)
			return
		}
		fmt.Println("User %s run command %s\n", usr.Name, cmd.Text)

	}

	var clicmd Cfcmd
	if clicmd.ConstructCmd(cmd.Text){

		err, ok := clicmd.RunCmd(&msg)
		if !ok{
			msg.Text = "Error executing command err: " + err.Error()
		}
	}else{
		msg.Text = "Bad command " + cmd.Text
	}



	msg.ResponseType = "in_channel"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)

}
	var slackApi *slack.Client


func main() {
	//retrieving the slack web api token from the environment variable
	access_token = os.Getenv("TOKEN")
	mongo_url := os.Getenv("MONGO")

	if mongo_url == "" {
		mongo_url = "mongodb://localhost:27017"
	} else {
		mongo_url = "mongodb://" + mongo_url + ":27017"
	}

	fmt.Printf("Welcome to Codefresh Slackbot\n")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := mongo.Connect(ctx, "mongodb://localhost:27017")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Mongo Connection Successful!\n")
	}

	if access_token == "" || access_token == "not_set" {
		fmt.Printf("WARNING: no access token set value is:%s\n", access_token)
	} else {
		fmt.Printf("Token set is:%s\n", access_token)

	}
	slackApi = slack.New(access_token)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handler)
	router.HandleFunc("/action", handleAction)
	router.HandleFunc("/action", handleAction)
	log.Fatal(http.ListenAndServe(":8080", router))

}

func handleAction(w http.ResponseWriter, r *http.Request) {


	var action slackActionMsg
	//var rsp slackRsp

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	//extracting the command
	err := action.ExecuteAction(r, w , true)

	if err != true {
		fmt.Printf("Cannot execute %s", r.Body)
		http.Error(w, "Cannot Parse", 400)
		return
	}




}