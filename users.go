package main



type User struct {
	Name	string `json:"name"`
	Token 	string `json:"token"`
}

var users = make(map[string]User)