package dotastats

import (
	"fmt"
	"strconv"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Mongodb struct {
	URI        string
	Dbname     string
	Collection string
}

func (mongo *Mongodb) SaveMatches(matchList []Match) error {
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return err
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	fmt.Println("saving matches")

	for _, match := range matchList {
		upsertdata := bson.M{"$set": match}
		condition := bson.M{"url": match.URL, "type": match.MatchType}
		info, err := collection.Upsert(condition, upsertdata)
		if err != nil {
			fmt.Errorf("error upserting %s", info, err)
		}
	}
	fmt.Println("done saving %v matches", len(matchList))
	return nil
}

func (mongo *Mongodb) GetTeamMatches(teamName, limit string) ([]Match, error) {
	var result []Match
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return []Match{}, err
	}
	var limitInt int
	var skipInt int
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			return []Match{}, err
		}
	}

	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	err = collection.Find(bson.M{
		"$or": []bson.M{
			bson.M{"teama": teamName},
			bson.M{"teamb": teamName},
		}}).Limit(limitInt).All(&result)

	if err != nil {
		return []Match{}, err
	}

	return result, nil
}

func (mongo *Mongodb) GetTeamF10kMatches(teamName, limit string) ([]Match, error) {
	var result []Match
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return []Match{}, err
	}
	var limitInt int
	var skipInt int
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			return []Match{}, err
		}
	}

	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	err = collection.Find(bson.M{
		"$or": []bson.M{
			bson.M{"teama": teamName},
			bson.M{"teamb": teamName},
		},
		"$or": []bson.M{
			bson.M{"scorea": 10},
			bson.M{"scoreb": 10},
		},
	}).Limit(limitInt).All(&result)

	if err != nil {
		return []Match{}, err
	}

	return result, nil
}
