package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Razielt77/cf-webapi-go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/nlopes/slack"
)



const (
	Mongo_DB 			= "slackbot"
	Mongo_users_collection	 = "users"
	Mongo_teams_collection	 = "teams"
)

type User struct {
	TeamID 			string `json:"teamid"`
	UserID			string `json:"userid"`
	Name			string `json:"name"`
	Team 			string `json:"team"`
	Avatar			string	`json:"avatar"`
	CFUserName		string `json:"cf_username"`
	ActiveAccount	string `json:"active_account"`
	CFAccounts []webapi.AccountInfo `json:"cf_accounts"`
}



func (u *User)SetToken(token string) error{

	fmt.Printf("Setting Token:\nToken is:%s\nActive Account is:%s\n",token,u.ActiveAccount)

	for i, _:= range u.CFAccounts{
		if u.CFAccounts[i].Name == u.ActiveAccount{
			u.CFAccounts[i].Token = token
			return nil
		}
	}
	return errors.New("no active account set")
}


func (u *User)GetToken() string{
	for _, account := range u.CFAccounts{
		if account.Name == u.ActiveAccount{
			return account.Token
		}
	}
	return ""
}

func (u *User)Print() {
	bytes, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Printf("User is:%s\n",bytes)
	}
}



func ensureIndex(s *mgo.Session) {

	session := s.Copy()
	defer session.Close()

	c := GetUsersCollection(session)

	index := mgo.Index{
		Key:        []string{"teamid","userid"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		fmt.Println(err)
	}
}




func AddUser(s *mgo.Session, user *User) (string, error) {

	session := s.Copy()
	defer session.Close()

	c := GetUsersCollection(session)

	err := c.Insert(user)
	if err != nil {
		if mgo.IsDup(err) {
			fmt.Printf("Failed add user: User with this ID:%s already exists\n",user.UserID)
			return "",err
		}

		fmt.Printf("Failed add user: Database error:%d\n",err)
		return "",err
	}

	fmt.Printf("User Added\n")

	return user.UserID,err
}

func GetUsersCollection(s *mgo.Session) *mgo.Collection {

	return s.DB(Mongo_DB).C(Mongo_users_collection)
}


func GetUser(s *mgo.Session,teamid string, userid string) (*User, error) {

	session := s.Copy()
	defer session.Close()

	c := GetUsersCollection(session)


	user := &User{}
	err := c.Find(bson.M{"teamid":teamid,"userid": userid}).One(user)

	if err != nil{
		if err == mgo.ErrNotFound{
			fmt.Printf("User not found\n")
		}else{
			fmt.Printf("Failed find user: Database error:%d\n",err)
		}
		return nil, err
	}

	return user,err

}


func UpdateUser(s *mgo.Session, user *User) error {

		session := s.Copy()
		defer session.Close()

		c := GetUsersCollection(session)

		err := c.Update(bson.M{"teamid": user.TeamID,"userid":user.UserID}, user)
		if err != nil {
			switch err {
			case mgo.ErrNotFound:
				fmt.Printf( "User not found. User Name: %s\n", user.Name)
				return err
			default:
				fmt.Printf( "Database error: %s\n", user.Name)
				return err
			}
		}

		return err

}

func SetUserToken (s *mgo.Session, callback *slack.InteractionCallback) bool {

	//user := User{TeamID:callback.Team.ID,UserID:callback.User.ID,Name:callback.User.Name,Team:callback.Team.Name}
	session := s.Copy()
	defer session.Close()

	user, _ := GetUser(session,callback.Team.ID,callback.User.ID)
	token := callback.Submission["cftoken"]

	if user == nil{


		user = &User{TeamID:callback.Team.ID,UserID:callback.User.ID,Name:callback.User.Name,Team:callback.Team.Name}

		//retrieving user's accounts

		cf_user, err := webapi.New(token).UserInfo()

		if err != nil {
			SendSimpleText(callback.ResponseURL,":heavy_exclamation_mark: *Invalid token*: "+ err.Error())
			return false
		}

		user.CFUserName = cf_user.Name
		user.CFAccounts = cf_user.Accounts
		user.ActiveAccount = cf_user.ActiveAccount
		user.Avatar = cf_user.UserData.Image
		err = user.SetToken(token)

		if err != nil{
			fmt.Println(err)
		}

		AddUser(session,user)

	}else{
		//fmt.Printf("Token submitted to exisiting account")

		//checking that tken is valid
		cf_user, err := webapi.New(token).UserInfo()

		if err != nil {
			SendSimpleText(callback.ResponseURL,":heavy_exclamation_mark: *Invalid token*: "+ err.Error())
			return false
		}

		if(cf_user.ActiveAccount != callback.State){
			SendSimpleText(callback.ResponseURL,":heavy_exclamation_mark: *Invalid token*: Token doesn't match with selected account")
			return false
		}

		user.ActiveAccount = callback.State
		user.SetToken(token)
		UpdateUser(s,user)
	}

	msg := slack.Msg{Text: ":white_check_mark: *Token successfully submitted!*"}
	att := slack.Attachment{
		Color:"#11b5a4",
		Text: "Welcome *"+user.CFUserName +
			"!*\nActive account is: *" +
			user.ActiveAccount +
			"*\nCurrently supported commands:\n" +
			"*/cf-pipelines-list*  Lists pipelines.\n"+
			"*/cf-pipelines-list-active*  Lists pipelines active past week.\n" +
			"*/cf-switch-account* Switch between your Codefresh's accounts.\n",
		ThumbURL: user.Avatar}

	msg.Attachments = append(msg.Attachments,att)

	DoPost(callback.ResponseURL,msg)


	return true
}



