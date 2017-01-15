package dotastats

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Match struct {
	Id         bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	TeamA      string        `json:"teama,omitempty" bson:"teama,omitempty"`
	TeamB      string        `json:"teamb,omitempty" bson:"teamb,omitempty"`
	MatchName  string        `json:"matchname,omitempty" bson:"matchname,omitempty"`
	URL        string        `json:"url,omitempty" bson:"url,omitempty"`
	Time       time.Time     `json:"time,omitempty" bson:"time,omitempty"`
	Tournament string        `json:"tournament,omitempty" bson:"tournament,omitempty"`
	MatchType  []string      `json:"type,omitempty" bson:"type,omitempty"`
	RatioA     float64       `json:"ratioa,omitempty" bson:"ratioa,omitempty"`
	RatioB     float64       `json:"ratiob,omitempty" bson:"ratiob,omitempty"`
	Note       string        `json:"note,omitempty" bson:"note,omitempty"`
	MatchID    int           `json:"matchid,omitempty" bson:"matchid,omitempty"`
	BestOf     string        `json:"bestof,omitempty" bson:"bestof,omitempty"`
	ScoreA     int           `json:"scorea" bson:"scorea"`
	ScoreB     int           `json:"scoreb" bson:"scoreb"`
	Winner     string        `json:"winner,omitempty" bson:"winner,omitempty"`
}
