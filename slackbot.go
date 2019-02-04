package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
)



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
		clicmd.RunCmd(&rsp)
	}else{
		rsp.Text = "Bad command " + cmd.Text
	}




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

	var cmd actionMsg
	var rsp slackRsp

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	//extracting the command
	err := cmd.extractMsg(r, true)

	if err != true {
		fmt.Println("Cannot parse %s", r.Body)
		http.Error(w, "Cannot Parse", 400)
		return
	}


	rsp.composeLoginScs(cmd)

	/*usr, ok := users[cmd.User_id]

	if !ok {
		rsp.composeLogin()
		//users[cmd.User_id] = User{Name:cmd.User_name}
	}else{
		rsp.Text = "User " + usr.Name + " exist."
	}*/


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rsp)

}