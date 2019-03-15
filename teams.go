package main

import (
	"fmt"
	"github.com/Razielt77/cf-webapi-go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Team struct {
	TeamID 			string `json:"teamid"`
	Team 			string `json:"team"`
	ActiveAccount	string `json:"active_account"`
	CFAccounts []webapi.AccountInfo `json:"cf_accounts"`
}


func ensureTeamIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := GetTeamsCollection(session)

	index := mgo.Index{
		Key:        []string{"teamid"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func GetTeamsCollection(s *mgo.Session) *mgo.Collection {

	return s.DB(Mongo_DB).C(Mongo_teams_collection)
}


func AddTeam(s *mgo.Session, team *Team) (string, error) {

	session := s.Copy()
	defer session.Close()

	c := GetTeamsCollection(session)

	err := c.Insert(team)
	if err != nil {
		if mgo.IsDup(err) {
			fmt.Printf("Failed add team: Team with this ID:%s and Name: %s already exists\n",team.TeamID,team.Team)
			return "",err
		}

		fmt.Printf("Failed add team: Database error:%d\n",err)
		return "",err
	}

	fmt.Printf("User Added\n")

	return team.TeamID,err
}

func GetTeam(s *mgo.Session,teamid string) (*Team, error) {

	session := s.Copy()
	defer session.Close()

	c := GetTeamsCollection(session)


	team := &Team{}
	err := c.Find(bson.M{"teamid":teamid}).One(team)

	if err != nil{
		if err == mgo.ErrNotFound{
			fmt.Printf("Team not found\n")
		}else{
			fmt.Printf("Failed find team: Database error:%d\n",err)
		}
		return nil, err
	}

	return team,err

}

func UpdateTeam(s *mgo.Session, team *Team) error {

	session := s.Copy()
	defer session.Close()

	c := GetTeamsCollection(session)

	err := c.Update(bson.M{"teamid": team.TeamID}, team)
	if err != nil {
		switch err {
		case mgo.ErrNotFound:
			fmt.Printf( "Team not found. Team Name: %s\n", team.Team)
			return err
		default:
			fmt.Printf( "Database error: %s\n", team.Team)
			return err
		}
	}

	return err

}