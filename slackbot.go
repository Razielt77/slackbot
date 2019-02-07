package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
)


var access_token = "xoxp-4027524866-4029668165-544352979810-1248b3b9011bf3035eb31be7380a144e"

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

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handler)
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
	err := action.ExtractAction(r, true)

	if err != true {
		fmt.Printf("Cannot parse %s", r.Body)
		http.Error(w, "Cannot Parse", 400)
		return
	}

	//executing the action
	err = action.ExecuteAction()

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
	json.NewEncoder(w).Encode(rsp)

}