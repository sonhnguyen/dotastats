package dotastats

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Match struct {
	Id         bson.ObjectId `json:"id" bson:"_id"`
	TeamA      string        `json:"teama" bson:"teama"`
	TeamB      string        `json:"teamb" bson:"teamb"`
	URL        string        `json:"url" bson:"url"`
	Time       time.Time     `json:"time" bson:"time"`
	Tournament string        `json:"tournament" bson:"tournament"`
	MatchType  []string      `json:"type" bson:"type"`
	RatioA     float64       `json:"ratioa" bson:"ratioa"`
	RatioB     float64       `json:"ratiob" bson:"ratiob"`
	Note       string        `json:"note" bson:"note"`
	MatchID    int           `json:"matchid" bson:"matchid"`
	BestOf     string        `json:"bestof" bson:"bestof"`
	ScoreA     int           `json:"scorea" bson:"scorea"`
	ScoreB     int           `json:"scoreb" bson:"scoreb"`
	Winner     string        `json:"winner" bson:"winner"`
}
