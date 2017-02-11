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
		condition := bson.M{"url": match.URL}
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
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			return []Match{}, err
		}
	}

	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	regexName := bson.M{"$regex": bson.RegEx{Pattern: "\\b" + teamName + "\\b", Options: "i"}}

	err = collection.Find(bson.M{
		"$or": []bson.M{
			bson.M{"teama": regexName},
			bson.M{"teamb": regexName},
			bson.M{"teama_short": regexName},
			bson.M{"teamb_short": regexName},
		}}).Limit(limitInt).Sort("-time").All(&result)

	if err != nil {
		return []Match{}, err
	}

	return result, nil
}

func (mongo *Mongodb) GetMatches(limit, skip, status string, fields []string) ([]Match, error) {
	var result []Match
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return []Match{}, err
	}
	var limitInt int
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			return []Match{}, err
		}
	}

	if skip != "" {
		skipInt, err = strconv.Atoi(skip)
		if err != nil {
			return []Match{}, err
		}
	}

	// Mapping to db format
	switch status {
	case "open":
		status = "Upcoming"
	case "closed":
		status = "Settled"
	case "live":
		status = "Live"
	default:
		status = "all"
	}

	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	if status != "all" {
		err = collection.Find(bson.M{"status": status}).Skip(skipInt).Limit(limitInt).Sort("-time").All(&result)
		if err != nil {
			return []Match{}, err
		}
	} else {
		var openMatches []Match
		err = collection.Find(bson.M{"status": "Upcoming"}).Skip(skipInt).Limit(limitInt).Sort("-time").All(&openMatches)
		if err != nil {
			return []Match{}, err
		}

		var liveMatches []Match
		err = collection.Find(bson.M{"status": "Live"}).Skip(skipInt).Limit(limitInt).Sort("-time").All(&liveMatches)
		if err != nil {
			return []Match{}, err
		}

		var closedMatches []Match
		err = collection.Find(bson.M{"status": "Settled"}).Skip(skipInt).Limit(limitInt).Sort("-time").All(&closedMatches)
		if err != nil {
			return []Match{}, err
		}
		result = append(closedMatches, openMatches...)
		result = append(result, liveMatches...)
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
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			return []Match{}, err
		}
	}

	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	regexName := bson.M{"$regex": bson.RegEx{Pattern: "\\b" + teamName + "\\b", Options: "i"}}
	regex10kills := bson.M{"$regex": bson.RegEx{Pattern: "10kills", Options: "i"}}

	err = collection.Find(bson.M{
		"$and": []bson.M{
			bson.M{"$or": []bson.M{
				bson.M{"teama": regexName},
				bson.M{"teamb": regexName},
				bson.M{"teama_short": regexName},
				bson.M{"teamb_short": regexName},
			}},
			bson.M{"mode_name": regex10kills},
			bson.M{"status": "Settled"},
		}},
	).Limit(limitInt).Sort("-time").All(&result)

	if err != nil {
		return []Match{}, err
	}

	return result, nil
}
