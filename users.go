package main

import (
	"fmt"
	"github.com/Razielt77/cf-webapi-go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)



const (
	Mongo_DB 			= "slackbot"
	Mongo_collection	 = "users"
)

type User struct {
	TeamID 		string `json:"teamid"`
	UserID		string `json:"userid"`
	Name		string `json:"name"`
	Team 		string `json:"team"`
	Token 		string `json:"token"`
	Avatar		string	`json:"avatar"`
	CFUserName	string `json:"cf_username"`
	DefaultAccount	int `json:"default_account"`
	CFAccounts []webapi.AccountInfo `json:"cf_accounts"`
}



func ensureIndex(s *mgo.Session) {

	session := s.Copy()
	defer session.Close()

	c := GetCollection(session)

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

	c := GetCollection(session)

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

func GetCollection(s *mgo.Session) *mgo.Collection {

	return s.DB(Mongo_DB).C(Mongo_collection)
}


func GetUser(s *mgo.Session,teamid string, userid string) (*User, error) {

	session := s.Copy()
	defer session.Close()

	c := GetCollection(session)


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

		c := GetCollection(session)

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


var users = make(map[string]User)