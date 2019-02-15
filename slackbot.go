package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
)

const NOT_FOUND = "not found"

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

	fmt.Printf("Welcome to Codefresh Slackbot\n")

	//retrieving the slack web api token from the environment variable
	access_token = os.Getenv("TOKEN")
	mongo_url := os.Getenv("MONGO")

	if mongo_url == "" {
		mongo_url = "localhost"
	}

	fmt.Printf("connecting o to %s\n",mongo_url)
	session, err := mgo.Dial(mongo_url)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	ensureIndex(session)

	/*user := User{TeamID:"2",UserID:"1",Name:"Raziel",Team:"Codefresh",CFTokens:[]CodefreshToken{{AccountName:`Codefresh-inc`, Token:`1111`},{AccountName:`Razielt77`,Token:`2222`}}}

	AddUser(session,&user)

	user2, err := GetUser(session,"2", "1")

	//fmt.Println(err)
	if user2 == nil{
		if err.Error() == NOT_FOUND {
			fmt.Printf("User not found\n")
		}else{
			fmt.Printf("Database Error: %s\n", err)
		}

	}else{
		fmt.Printf("User Found\nName: %v\n",user2.Name)
		for _, s := range user2.CFTokens{
			fmt.Printf("Account Name: %v\n",s.AccountName)
		}
	}*/



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