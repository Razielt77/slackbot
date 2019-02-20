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

	db := session.DB(Mongo_DB)
	err = db.DropDatabase()
	if err != nil {
		fmt.Printf("cannot drop %s\n",err)
	}

	ensureIndex(session)

	/*user := User{TeamID:"2",UserID:"1",Name:"Raziel",Team:"Codefresh",CFTokens:[]CodefreshToken{{AccountName:`Codefresh-inc`, Token:`1111`},{AccountName:`Razielt77`,Token:`2222`}}}

	AddUser(session,&user)

	user2, err := GetUser(session,"2", "1")

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
	}

	user2.CFTokens = append(user.CFTokens,CodefreshToken{AccountName:"Dustin",Token:"4444"})

	UpdateUser(session,user2)


	user3, err := GetUser(session,"2", "1")

	if user3 == nil{
		if err.Error() == NOT_FOUND {
			fmt.Printf("User not found\n")
		}else{
			fmt.Printf("Database Error: %s\n", err)
		}

	}else{
		fmt.Printf("User Found\nName: %v\n",user3.Name)
		for _, s := range user3.CFTokens{
			fmt.Printf("Account Name: %v\n",s.AccountName)
		}
	}

	if access_token == "" || access_token == "not_set" {
		fmt.Printf("WARNING: no access token set value is:%s\n", access_token)
	} else {
		fmt.Printf("Token set is:%s\n", access_token)

	}*/


	slackApi = slack.New(access_token)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Handler(session))
	router.HandleFunc("/action", HandleAction(session))
	router.HandleFunc("/pipelinelist", PipelineListAction(session))
	log.Fatal(http.ListenAndServe(":8080", router))

}


func HandleAction (s *mgo.Session) func(w http.ResponseWriter, r *http.Request){
	return func (w http.ResponseWriter, r *http.Request){

		session := s.Copy()
		defer session.Close()

		var action slackActionMsg
		//var rsp slackRsp

		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		//extracting the command
		err := action.ExecuteAction(session,r, w , true)

		if err != true {
			fmt.Printf("Cannot execute %s", r.Body)
			http.Error(w, "Cannot Parse", 400)
			return
		}

	}
}



func PipelineListAction (s *mgo.Session) func(w http.ResponseWriter, r *http.Request){
	return func (w http.ResponseWriter, r *http.Request){

		session := s.Copy()
		defer session.Close()

		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		cmd , err := slack.SlashCommandParse(r)

		if err != nil {
			fmt.Println("Cannot parse %s\n", r.Body)
			http.Error(w, "Cannot Parse", 400)
			return
		}

		response_url := cmd.ResponseURL

		usr, err := GetUser(session,cmd.TeamID,cmd.UserID)

		msg := slack.Msg{}
		msg.ResponseType = "in_channel"
		w.Header().Set("Content-Type", "application/json")


		if usr == nil {
			composeLogin(&msg)
			json.NewEncoder(w).Encode(msg)
			return
		}

		//fmt.Printf("User found: %s\nUsing token:%s\n",usr.Name,usr.CFTokens[0].Token)

		msg.Text = "Retrieving  Pipelines..."
		json.NewEncoder(w).Encode(msg)


		go SendPipelinesListMsg(usr,response_url)


		//json.NewEncoder(w).Encode(msg)
		return

	}
}



func Handler(s *mgo.Session) func(w http.ResponseWriter, r *http.Request){
	return func (w http.ResponseWriter, r *http.Request){

		session := s.Copy()
		defer session.Close()

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
			usr, _ := GetUser(session,cmd.Team_id,cmd.User_id)
			if usr == nil{
				composeLogin(&msg)
				//users[cmd.User_id] = User{Name:cmd.User_name}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(msg)
				return
			}
			fmt.Printf("User %s run command %s\n", usr.Name, cmd.Text)
		}

		/*var clicmd Cfcmd
		if clicmd.ConstructCmd(cmd.Text){

			err, ok := clicmd.RunCmd(&msg)
			if !ok{
				msg.Text = "Error executing command err: " + err.Error()
			}
		}else{
			msg.Text = "Bad command " + cmd.Text
		}*/



		msg.ResponseType = "in_channel"

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)

	}
}
