package dotastats

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Mongodb struct {
	URI            string
	Dbname         string
	Collection     string
	CollectionTeam string
}

func selectFields(q ...string) (r bson.M) {
	r = make(bson.M, len(q))
	for _, s := range q {
		r[s] = 1
	}
	return
}

func buildTeamInfoQuery(apiParams APIParams) bson.M {
	r := make(bson.M, 1)

	if apiParams.Game != "" && apiParams.Game != "all" {
		r["game"] = apiParams.Game
	}

	return r
}
func buildFindQuery(apiParams APIParams) bson.M {
	r := make(bson.M, 2)

	if apiParams.Game != "" && apiParams.Game != "all" {
		r["game"] = apiParams.Game
	}
	r["time"] = bson.M{"$gt": apiParams.TimeFrom,
		"$lt": apiParams.TimeTo,
	}

	return r
}

func (mongo *Mongodb) SaveTeamInfo(teamList []TeamInfo) error {
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return err
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB(mongo.Dbname).C(mongo.CollectionTeam)
	fmt.Println("saving teamInfo")

	for _, team := range teamList {

		upsertdata := bson.M{"$set": team}
		condition := bson.M{"url": team.URL}
		info, err := collection.Upsert(condition, upsertdata)
		if err != nil {
			fmt.Errorf("error upserting %s", info, err)
		}
	}
	fmt.Println("done saving %v teams", len(teamList))
	return nil
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

func (mongo *Mongodb) GetAllTeamInfo() ([]TeamInfo, error) {
	var result []TeamInfo
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return []TeamInfo{}, err
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB(mongo.Dbname).C(mongo.CollectionTeam)
	err = collection.Find(bson.M{}).All(&result)

	if err != nil {
		return []TeamInfo{}, err
	}

	return result, nil
}

func (mongo *Mongodb) GetTeamInfo(teamSlug string, apiParams APIParams) (TeamInfo, error) {
	var result TeamInfo
	var findQuery bson.M
	findQuery = buildTeamInfoQuery(apiParams)
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return TeamInfo{}, err
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB(mongo.Dbname).C(mongo.CollectionTeam)
	regexSlug := bson.M{"$regex": bson.RegEx{Pattern: "\\b" + teamSlug + "\\b", Options: "i"}}
	findQuery["slug"] = regexSlug
	err = collection.Find(findQuery).One(&result)

	if err != nil {
		return TeamInfo{}, err
	}

	return result, nil
}

func (mongo *Mongodb) GetTeamHistoryMatches(teamA, teamB string, apiParams APIParams) ([]Match, error) {
	var result []Match
	var findQuery bson.M
	findQuery = buildFindQuery(apiParams)
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return []Match{}, err
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	regexMode := bson.M{"$regex": bson.RegEx{Pattern: "(\\bMatch Winner\\b|\\b10kills\\b)", Options: "i"}}
	findQuery["$and"] = []bson.M{
		bson.M{"$or": []bson.M{
			bson.M{"$and": []bson.M{
				bson.M{"teama": teamA},
				bson.M{"teamb": teamB},
			}},
			bson.M{"$and": []bson.M{
				bson.M{"teama": teamB},
				bson.M{"teamb": teamA},
			}},
		}},
		bson.M{"status": "Settled"},
		bson.M{"mode_name": regexMode},
	}

	err = collection.Find(findQuery).Select(selectFields(apiParams.Fields...)).Skip(apiParams.Skip).Limit(apiParams.Limit).Sort("-time").All(&result)

	if err != nil {
		return []Match{}, err
	}
	return result, nil
}

func (mongo *Mongodb) GetTeamMatches(teamName string, apiParams APIParams) ([]Match, error) {
	var result []Match
	var findQuery bson.M
	findQuery = buildFindQuery(apiParams)
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return []Match{}, err
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	regexMode := bson.M{"$regex": bson.RegEx{Pattern: "(\\bMatch Winner\\b|\\b10kills\\b)", Options: "i"}}
	findQuery["$and"] = []bson.M{
		bson.M{"$or": []bson.M{
			bson.M{"teama": teamName},
			bson.M{"teamb": teamName},
			bson.M{"teama_short": teamName},
			bson.M{"teamb_short": teamName},
		}},
		bson.M{"status": "Settled"},
		bson.M{"mode_name": regexMode},
	}
	err = collection.Find(findQuery).Select(selectFields(apiParams.Fields...)).Skip(apiParams.Skip).Limit(apiParams.Limit).Sort("-time").All(&result)

	if err != nil {
		return []Match{}, err
	}
	return result, nil
}

func (mongo *Mongodb) GetMatchesList(status string, apiParams APIParams) ([]Match, error) {
	var result []Match
	var findQuery bson.M
	findQuery = buildFindQuery(apiParams)
	sess, err := mgo.Dial(mongo.URI)
	fmt.Println(mongo.URI)
	if err != nil {
		return []Match{}, err
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	// Mapping to db format
	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	if status == "all" {
		var openMatches []Match
		findQuery["status"] = "Upcoming"
		err = collection.Find(findQuery).Select(selectFields(apiParams.Fields...)).Skip(apiParams.Skip).Limit(apiParams.Limit).Sort("time").All(&openMatches)
		if err != nil {
			return []Match{}, err
		}

		findQuery["status"] = "Live"
		var liveMatches []Match
		err = collection.Find(findQuery).Select(selectFields(apiParams.Fields...)).Skip(apiParams.Skip).Limit(apiParams.Limit).Sort("-time").All(&liveMatches)
		if err != nil {
			return []Match{}, err
		}

		findQuery["status"] = "Settled"
		var closedMatches []Match
		err = collection.Find(findQuery).Select(selectFields(apiParams.Fields...)).Skip(apiParams.Skip).Limit(apiParams.Limit).Sort("-time").All(&closedMatches)
		if err != nil {
			return []Match{}, err
		}
		result = append(liveMatches, openMatches...)
		result = append(result, closedMatches...)
	} else {
		sort := "-time"
		switch status {
		case "open":
			findQuery["status"] = "Upcoming"
			sort = "time"
		case "closed":
			findQuery["status"] = "Settled"
		case "live":
			findQuery["status"] = "Live"
		default:
			return []Match{}, err
		}
		err = collection.Find(findQuery).Select(selectFields(apiParams.Fields...)).Skip(apiParams.Skip).Limit(apiParams.Limit).Sort(sort).All(&result)
		if err != nil {
			return []Match{}, err
		}
	}
	return result, nil
}

func (mongo *Mongodb) GetMatchByID(matchID string) (Match, error) {
	var result Match
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return Match{}, err
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	if bson.IsObjectIdHex(matchID) {
		err = collection.FindId(bson.ObjectIdHex(matchID)).One(&result)
		if err != nil {
			return Match{}, err
		}
	} else {
		return Match{}, fmt.Errorf("Invalid input in ID %s", matchID)
	}
	return result, nil
}

func (mongo *Mongodb) GetTeamF10kMatches(teamName string, apiParams APIParams) ([]Match, error) {
	var result []Match
	findQuery := buildFindQuery(apiParams)
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return []Match{}, err
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	regex10kills := bson.M{"$regex": bson.RegEx{Pattern: "10kills", Options: "i"}}
	findQuery["$and"] = []bson.M{
		bson.M{"$or": []bson.M{
			bson.M{"teama": teamName},
			bson.M{"teamb": teamName},
			bson.M{"teama_short": teamName},
			bson.M{"teamb_short": teamName},
		}},
		bson.M{"mode_name": regex10kills},
		bson.M{"status": "Settled"},
	}
	err = collection.Find(findQuery).Select(selectFields(apiParams.Fields...)).Skip(apiParams.Skip).Limit(apiParams.Limit).Sort("-time").All(&result)

	if err != nil {
		return []Match{}, err
	}
	return result, nil
}
