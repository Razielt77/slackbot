package main



type User struct {
	Name	string `json:"name"`
	Token 	string `json:"token"`
	CfToken string `json:"cf_token"`
}



var users = make(map[string]User)