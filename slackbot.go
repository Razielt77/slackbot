package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
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
	//var slackApi *slack.Client


func main() {
	//retrieving the slack web api token from the environment variable
	access_token = os.Getenv("TOKEN")

	api := slack.New("YOUR TOKEN HERE")
	//logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	//slack.SetLogger(logger)
	//api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace C2147483705 with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "C2147483705"))

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}

	/*if access_token == "" || access_token == "not_set"{
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
	//http.ListenAndServe(":8080", nil)*/
}


func handleAction(w http.ResponseWriter, r *http.Request) {


	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()

	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{access_token}))
	if e != nil {
		fmt.Printf("Error Parsin %s\n", r.Body)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Printf("eventsAPIEvent is %s\n", eventsAPIEvent)



	var action slackActionMsg
	var rsp slackRsp

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	//extracting the command
	err := action.ExtractAction(r, true)

	if err != true {
		fmt.Printf("Cannot parse %s", r.Body)
		http.Error(w, "Cannot Parse", 400)
		return
	}


	switch action.Type {
	case "interactive_message":
		//executing the action
		err = action.ExecuteAction()
	case "dialog_submission":
		err = action.DialogSubmission()

	}

	if err != true {
		fmt.Printf("Cannot execute %v", action.Actions[0])
		http.Error(w, "Cannot Parse", 400)
		return
	}


	//rsp.composeLoginScs(cmd)

	/*usr, ok := users[cmd.User_id]

	if !ok {
		rsp.composeLogin()
		//users[cmd.User_id] = User{Name:cmd.User_name}
	}else{
		rsp.Text = "User " + usr.Name + " exist."
	}*/

	rsp.ResponseType = "in_channel"
	rsp.Text = "In Action baby"


	w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(rsp)

}