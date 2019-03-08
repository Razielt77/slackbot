package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
)

const (
	NOT_FOUND = "not found"
	SLACK_TOKEN_ENV_NAME = "TOKEN"
	MONGO_URL_RNV_NAME = "MONGO"
)

var access_token string = ""





var slackApi *slack.Client


func main() {

	fmt.Printf("Welcome to Codefresh Slackbot\n")

	//retrieving the slack web api token from the environment variable
	access_token = os.Getenv(SLACK_TOKEN_ENV_NAME)
	mongo_url := os.Getenv(MONGO_URL_RNV_NAME)

	if mongo_url == "" {
		mongo_url = "localhost"
	}

	session, err := mgo.Dial(mongo_url)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	/*db := session.DB(Mongo_DB)
	err = db.DropDatabase()
	if err != nil {
		fmt.Printf("cannot drop %s\n",err)
	}*/

	ensureIndex(session)

	slackApi = slack.New(access_token)


	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Handler(session))
	router.HandleFunc("/action", HandleAction(session))
	router.HandleFunc("/accountchange", AccountChangeCommand(session))
	router.HandleFunc("/pipelineslist", PipelineListAction(session))
	router.HandleFunc("/events", HandleEvent(session))
	log.Fatal(http.ListenAndServe(":8080", router))

}


func HandleAction (s *mgo.Session) func(w http.ResponseWriter, r *http.Request){
	return func (w http.ResponseWriter, r *http.Request){

		w.WriteHeader(http.StatusOK)
		session := s.Copy()
		defer session.Close()

		var action slackActionMsg
		//var rsp slackRsp

		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		//extracting the command
		action.ExecuteAction(s,r, true)


	}
}

func HandleEvent (s *mgo.Session) func(w http.ResponseWriter, r *http.Request){
	return func (w http.ResponseWriter, r *http.Request){

		w.WriteHeader(http.StatusOK)
		session := s.Copy()
		defer session.Close()



		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}else{
			fmt.Printf("Body: %s\n",r.Body)
		}


		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: "TOKEN"}))
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				slackApi.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			}
		}



	}
}


func PipelineListAction (s *mgo.Session) func(w http.ResponseWriter, r *http.Request){
	return func (w http.ResponseWriter, r *http.Request){

		w.WriteHeader(http.StatusOK)

		cmd := ParseSlashCommand(w,r)
		if cmd == nil {
			return
		}

		session := s.Copy()
		defer session.Close()

		usr, _ := GetUser(session,cmd.TeamID,cmd.UserID)

		if usr == nil {
			lgn := ComposeLogin()
			go DoPost(cmd.ResponseURL,lgn)
			return
		}

		go SendPipelinesListMsg(usr,cmd)

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

		var msg *slack.Msg
		//if command require login than check if user logged in (have a context) if not it asks him/her to login
		if cmd.LoginRequired() {
			usr, _ := GetUser(session,cmd.Team_id,cmd.User_id)
			if usr == nil{
				msg = ComposeLogin()
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(msg)
				return
			}
			fmt.Printf("User %s run command %s\n", usr.Name, cmd.Text)
		}




		msg.ResponseType = "in_channel"

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)

	}
}
