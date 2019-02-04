package main



type User struct {
	Name	string `json:"name"`
	Token 	string `json:"token"`
	CfContext string `json:"cfcontext"`
}



var users = make(map[string]User)