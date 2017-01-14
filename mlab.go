package dotastats

import (
	"fmt"

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
