package main

import (
	"net/http"
	"encoding/json"
)


type slackCmd struct {
	token    		string       //  `gIkuvaNzQIHg97ATvDxqgjtO`
	team_id  		string       // `T0001`
	team_domain   	string       // `example`
	enterprise_id	string 		//  'E0001'
	enterprise_name	string		//  Globular%20Construct%20Inc
	channel_id		string 		//	C2147483705
	channel_name	string	    //  test
	user_id			string 		// U2147483697
	user_name		string 		// Steve
	command			string 	  	//weather
	text			string		//94070
	response_url	string		//https://hooks.slack.com/commands/1234/5678

}


type slackRsp struct {
	text			string
}



func handler(w http.ResponseWriter, r *http.Request) {
	//var cmd slackCmd

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	/*err := json.NewDecoder(r.Body).Decode(&cmd)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}*/

	rsp := slackRsp{text: "All fine"}

	json.NewEncoder(w).Encode(rsp)

	//fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
