package main



type User struct {
	ID	string `json:"id"`
	Name	string `json:"name"`
	Team 	string `json:"team"`
	CFTokens []CodefreshToken `json:"cftokens"`
}

type CodefreshToken struct {
	AccountName 	string `json:"accountname"`
	Token 			string 	`json:"token"`
}



var users = make(map[string]User)