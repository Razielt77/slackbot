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
	MONGO_URL_ENV_NAME = "MONGO"
	SLACK_VER_ENV_NAME = "VER_TOKEN"
	ENTER_TOKEN = "enter_token"
	SWITCH_ACCOUNT = "switch_account"
	NOT_AVAILABLE string  = "Not Available"
	ACTIVE_PIPELINE_COMMAND string  = "/cf-pipelines-list-active"
	ACTIVE_DURATION_IN_HOURS float64  = 72
	PIPELINE_ACTION = "pipeline_action"
	VIEW_BUILDS = "view_builds"
	IN_CHANNEL = "in_channel"
	FIRST_TIME_USER = "It seems like your team never added a Codefresh's token. Please add your account's token."
)


var access_token string = ""
var ver_token string = ""





var slackApi *slack.Client


func main() {

	fmt.Printf("Welcome to Codefresh Slackbot\n")

	//retrieving the slack web api token from the environment variable
	access_token = os.Getenv(SLACK_TOKEN_ENV_NAME)
	ver_token = os.Getenv(SLACK_VER_ENV_NAME)
	mongo_url := os.Getenv(MONGO_URL_ENV_NAME)

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
	ensureTeamIndex(session)

	slackApi = slack.New(access_token)


	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Handler(session))
	router.HandleFunc("/action", HandleAction(session))
	router.HandleFunc("/tokenlist", ListTokens(session))
	router.HandleFunc("/accountchange", AccountChangeCommand(session))
	router.HandleFunc("/pipelineslist", PipelineListAction(session))
	router.HandleFunc("/events", HandleEvent(session))
	log.Fatal(http.ListenAndServe(":8080", router))

}




func HandleEvent (s *mgo.Session) func(w http.ResponseWriter, r *http.Request){
	return func (w http.ResponseWriter, r *http.Request){

		session := s.Copy()
		defer session.Close()



		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}


		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()

		//fmt.Printf("Body: %s\n",body)

		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: ver_token}))
		//eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body))

		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("Parsing error: %s\n",e.Error())
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
			case *slackevents.LinkSharedEvent:
				go EnrichSharedLink(s,eventsAPIEvent.TeamID,ev)
				fmt.Printf("here is the link %s\n",ev.Links[0].URL)
				w.WriteHeader(http.StatusOK)
				//slackApi.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
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
			lgn := ComposeLogin(FIRST_TIME_USER)
			go DoPost(cmd.ResponseURL,lgn)
			return
		}

		go SendPipelinesListMsg(usr,cmd)

		return

	}
}

func ListTokens (s *mgo.Session) func(w http.ResponseWriter, r *http.Request){
	return func (w http.ResponseWriter, r *http.Request){

		w.WriteHeader(http.StatusOK)

		cmd := ParseSlashCommand(w,r)
		if cmd == nil {
			return
		}

		session := s.Copy()
		defer session.Close()

		team, _ := GetTeam(session,cmd.TeamID)

		if team == nil {
			lgn := ComposeLogin(FIRST_TIME_USER)
			go DoPost(cmd.ResponseURL,lgn)
			return
		}

		go SendTokensList(team,cmd)

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
				msg = ComposeLogin(FIRST_TIME_USER)
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
