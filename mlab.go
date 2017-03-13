package dotastats

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Mongodb struct {
	URI        string
	Dbname     string
	Collection string
}

func selectFields(q ...string) (r bson.M) {
	r = make(bson.M, len(q))
	for _, s := range q {
		r[s] = 1
	}
	return
}
func filterGame(matches []Match, game string) []Match {
	var gameSelected string
	var result []Match
	gameList := []string{"dota", "csgo", "snooker", "football", "basketball"}
	for _, v := range gameList {
		if game == v {
			gameSelected = game
		}
	}
	if gameSelected != "" {
		for _, v := range matches {
			if v.Game == gameSelected {
				result = append(result, v)
			}
		}
		return result
	}
	return matches
}

func filterTime(matches []Match, timeFrom, timeTo time.Time) []Match {
	var result []Match
	for _, v := range matches {
		if v.Time.After(timeFrom) && v.Time.Before(timeTo) {
			result = append(result, v)
		}
	}
	return result
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

func (mongo *Mongodb) GetTeamMatches(teamName string, apiParams APIParams) ([]Match, error) {
	var result []Match
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return []Match{}, err
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	regexName := bson.M{"$regex": bson.RegEx{Pattern: "\\b" + teamName + "\\b", Options: "i"}}

	err = collection.Find(bson.M{
		"$or": []bson.M{
			bson.M{"teama": regexName},
			bson.M{"teamb": regexName},
			bson.M{"teama_short": regexName},
			bson.M{"teamb_short": regexName},
		}}).Select(selectFields(apiParams.Fields...)).Skip(apiParams.Skip).Limit(apiParams.Limit).Sort("-time").All(&result)

	if err != nil {
		return []Match{}, err
	}
	result = filterGame(result, apiParams.Game)
	result = filterTime(result, apiParams.TimeFrom, apiParams.TimeTo)
	return result, nil
}

func (mongo *Mongodb) GetMatches(status string, apiParams APIParams) ([]Match, error) {
	var result []Match
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return []Match{}, err
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

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
		err = collection.Find(bson.M{"status": status}).Select(selectFields(apiParams.Fields...)).Skip(apiParams.Skip).Limit(apiParams.Limit).Sort("-time").All(&result)
		if err != nil {
			return []Match{}, err
		}
	} else {
		var openMatches []Match
		err = collection.Find(bson.M{"status": "Upcoming"}).Select(selectFields(apiParams.Fields...)).Skip(apiParams.Skip).Limit(apiParams.Limit).Sort("-time").All(&openMatches)
		if err != nil {
			return []Match{}, err
		}

		var liveMatches []Match
		err = collection.Find(bson.M{"status": "Live"}).Select(selectFields(apiParams.Fields...)).Skip(apiParams.Skip).Limit(apiParams.Limit).Sort("-time").All(&liveMatches)
		if err != nil {
			return []Match{}, err
		}

		var closedMatches []Match
		err = collection.Find(bson.M{"status": "Settled"}).Select(selectFields(apiParams.Fields...)).Skip(apiParams.Skip).Limit(apiParams.Limit).Sort("-time").All(&closedMatches)
		if err != nil {
			return []Match{}, err
		}
		result = append(closedMatches, openMatches...)
		result = append(result, liveMatches...)
	}
	result = filterGame(result, apiParams.Game)
	result = filterTime(result, apiParams.TimeFrom, apiParams.TimeTo)
	return result, nil
}

func (mongo *Mongodb) GetTeamF10kMatches(teamName string, apiParams APIParams) ([]Match, error) {
	var result []Match
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return []Match{}, err
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

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
	).Select(selectFields(apiParams.Fields...)).Skip(apiParams.Skip).Limit(apiParams.Limit).Sort("-time").All(&result)

	if err != nil {
		return []Match{}, err
	}
	result = filterGame(result, apiParams.Game)
	result = filterTime(result, apiParams.TimeFrom, apiParams.TimeTo)
	return result, nil
}
