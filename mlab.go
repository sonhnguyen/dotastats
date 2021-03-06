package dotastats

import (
	"fmt"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Mongodb struct {
	URI                string
	Dbname             string
	Collection         string
	CollectionTeam     string
	CollectionDotaTeam string
	CollectionProMatch string
	CollectionFeedback string
	CollectionUser     string
	CollectionSession  string
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

func (mongo *Mongodb) SaveOpenDotaTeam(teamList []OpenDotaTeam) error {
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return err
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB(mongo.Dbname).C(mongo.CollectionDotaTeam)
	fmt.Println("saving CollectionDotaTeam")

	for _, team := range teamList {
		upsertdata := bson.M{"$set": team}
		condition := bson.M{"team_id": team.TeamID}
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

func (mongo *Mongodb) SaveOpenDotaProMatches(matchList []OpenDotaMatch) error {
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return err
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB(mongo.Dbname).C(mongo.CollectionProMatch)
	fmt.Println("saving matches")

	for _, match := range matchList {
		upsertdata := bson.M{"$set": match}
		condition := bson.M{"match_id": match.MatchID}
		info, err := collection.Upsert(condition, upsertdata)
		if err != nil {
			fmt.Errorf("error upserting %s", info, err)
		}
	}
	fmt.Println("done saving %v matches", len(matchList))
	return nil
}

func (mongo *Mongodb) GetFeedback() ([]Feedback, error) {
	var feedBackArr []Feedback
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return []Feedback{}, err
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB(mongo.Dbname).C(mongo.CollectionFeedback)

	err = collection.Find(nil).All(&feedBackArr)
	if err != nil {
		return []Feedback{}, fmt.Errorf("error getting feedback %s", err)
	}
	return feedBackArr, nil
}

func (mongo *Mongodb) SaveFeedback(feedBack *Feedback) error {
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return err
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB(mongo.Dbname).C(mongo.CollectionFeedback)

	err = collection.Insert(feedBack)
	if err != nil {
		fmt.Errorf("error inserting %s", err)
	}
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

func regexIncludeWord(s string) bson.M {
	pattern := fmt.Sprintf("(\\b%s\\b)", s)
	result := bson.M{"$regex": bson.RegEx{Pattern: pattern, Options: "i"}}
	return result
}

func (mongo *Mongodb) GetOpenDotaMatch(match Match) (OpenDotaMatch, bool, error) {
	var result OpenDotaMatch
	teamAIsRadiant := false
	findQuery := make(bson.M, 2)
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return OpenDotaMatch{}, false, err
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB(mongo.Dbname).C(mongo.CollectionProMatch)
	findQuery["$or"] = []bson.M{
		bson.M{"$and": []bson.M{
			bson.M{"radiant_name": regexIncludeWord(match.TeamA)},
			bson.M{"dire_name": regexIncludeWord(match.TeamB)},
		}},
		bson.M{"$and": []bson.M{
			bson.M{"radiant_name": regexIncludeWord(match.TeamB)},
			bson.M{"dire_name": regexIncludeWord(match.TeamA)},
		}},
		bson.M{"$and": []bson.M{
			bson.M{"radiant_tag": regexIncludeWord(match.TeamAShort)},
			bson.M{"dire_tag": regexIncludeWord(match.TeamBShort)},
		}},
		bson.M{"$and": []bson.M{
			bson.M{"radiant_tag": regexIncludeWord(match.TeamBShort)},
			bson.M{"dire_tag": regexIncludeWord(match.TeamAShort)},
		}},
	}

	findQuery["start_time"] = bson.M{"$lte": match.Time}
	err = collection.Find(findQuery).Sort("-start_time").Limit(1).One(&result)
	if err != nil {
		return OpenDotaMatch{}, false, err
	}

	if strings.Contains(result.RadiantName, match.TeamA) ||
		strings.Contains(result.RadiantName, match.TeamAShort) ||
		strings.Contains(result.RadiantTag, match.TeamA) ||
		strings.Contains(result.RadiantTag, match.TeamAShort) {
		teamAIsRadiant = true
	} else {
		teamAIsRadiant = false
	}

	return result, teamAIsRadiant, nil
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

func (mongo *Mongodb) GetTeamFBMatches(teamName string, apiParams APIParams) ([]Match, error) {
	var result []Match
	findQuery := buildFindQuery(apiParams)
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return []Match{}, err
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB(mongo.Dbname).C(mongo.Collection)
	regexfb := bson.M{"$regex": bson.RegEx{Pattern: "1st Blood", Options: "i"}}
	findQuery["$and"] = []bson.M{
		bson.M{"$or": []bson.M{
			bson.M{"teama": teamName},
			bson.M{"teamb": teamName},
			bson.M{"teama_short": teamName},
			bson.M{"teamb_short": teamName},
		}},
		bson.M{"mode_name": regexfb},
		bson.M{"status": "Settled"},
	}
	err = collection.Find(findQuery).Select(selectFields(apiParams.Fields...)).Skip(apiParams.Skip).Limit(apiParams.Limit).Sort("-time").All(&result)

	if err != nil {
		return []Match{}, err
	}
	return result, nil
}

func (mongo *Mongodb) GetUserByEmail(email string) (User, error) {
	var user User
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return User{}, err
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB(mongo.Dbname).C(mongo.CollectionUser)
	err = collection.Find(bson.M{"email": email}).One(&user)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (mongo *Mongodb) CreateUser(user *User) error {
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return err
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB(mongo.Dbname).C(mongo.CollectionUser)

	err = collection.Insert(user)
	if err != nil {
		return err
	}
	return nil
}

func (mongo *Mongodb) GetSessionBySessionKey(ssk string) (Session, error) {
	var dotastatsSession Session
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return Session{}, err
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB(mongo.Dbname).C(mongo.CollectionSession)
	err = collection.Find(bson.M{"session_key": ssk}).One(&dotastatsSession)

	if err != nil {
		return Session{}, err
	}

	return dotastatsSession, nil
}

func (mongo *Mongodb) CreateOrUpdateSession(s Session) error {
	sess, err := mgo.Dial(mongo.URI)
	if err != nil {
		return err
	}

	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB(mongo.Dbname).C(mongo.CollectionSession)
	count, err := collection.Find(bson.M{"email": s.Email}).Count()
	if err != nil {
		return err
	}

	if count > 0 {
		err = collection.Update(bson.M{"email": s.Email}, bson.M{"$set": bson.M{"session_key": s.SessionKey}})
		return err
	}

	err = collection.Insert(s)

	if err != nil {
		return err
	}

	return nil
}
