package main

import (
	"fmt"
	"net/http"
	"encoding/json"
)


type slackCmd struct {
	Token    		string      `json:"token"`
	Team_id  		string      `json:"team_id"`
	Team_domain   	string      `json:"team_domain"`
	Enterprise_id	string 		`json:"enterprise_id"`
	Enterprise_name	string		`json:"enterprise_name"`
	Channel_id		string 		`json:"channel_id"`
	Channel_name	string	    `json:"channel_name"`
	User_id			string 		`json:"user_id"`
	User_name		string 		`json:"user_name"`
	Command			string 	  	`json:"command"`
	Rext			string		`json:"text"`
	Response_url	string		`json:"response_url"`

}



type slackRsp struct {
	Text	string `json:"text"`
}



func handler(w http.ResponseWriter, r *http.Request) {

	//var cmd slackCmd

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	fmt.Println("Body- %s", r.Body)

	/*err := json.NewDecoder(r.Body).Decode(&cmd)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}*/


	var rsp slackRsp
	rsp.Text = "All fine"
	//rsp.Id = "11"

	/*js, err := json.Marshal(rsp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}*/


	json.NewEncoder(w).Encode(rsp)

	//fmt.Fprintf(w, "Hi there, I love %s! and %s", rsp, js)
	//fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)

	/*var rsp slackRsp
	rsp.text = "All fine"
	rsp.id = "dd"
	js, err := json.Marshal(rsp)

	if err != nil {
		fmt.Println("Error")
		return
	}
	fmt.Println("js - %s rsp - %s", js, rsp)*/


	http.ListenAndServe(":8080", nil)
}
