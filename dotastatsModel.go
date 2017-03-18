package dotastats

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type APIParams struct {
	Limit    int
	Skip     int
	Fields   []string
	TimeFrom time.Time
	TimeTo   time.Time
	Game     string
}

type Match struct {
	Id             bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	TeamAID        string        `json:"teama_id,omitempty" bson:"teama_id,omitempty"`
	TeamBID        string        `json:"teamb_id,omitempty" bson:"teamb_id,omitempty"`
	TeamA          string        `json:"teama,omitempty" bson:"teama,omitempty"`
	TeamB          string        `json:"teamb,omitempty" bson:"teamb,omitempty"`
	LogoA          string        `json:"logo_a,omitempty" bson:"logo_a,omitempty"`
	LogoB          string        `json:"logo_b,omitempty" bson:"logo_b,omitempty"`
	TeamAShort     string        `json:"teama_short,omitempty" bson:"teama_short,omitempty"`
	TeamBShort     string        `json:"teamb_short,omitempty" bson:"teamb_short,omitempty"`
	Tournament     string        `json:"tournament,omitempty" bson:"tournament,omitempty"`
	TournamentLogo string        `json:"tournament_logo,omitempty" bson:"tournament_logo,omitempty"`
	Game           string        `json:"game,omitempty" bson:"game,omitempty"`
	BestOf         string        `json:"bestof,omitempty" bson:"bestof,omitempty"`
	// sub match specific
	MatchID        string     `json:"matchid,omitempty" bson:"matchid,omitempty"`
	URL            string     `json:"url,omitempty" bson:"url,omitempty"`
	Time           *time.Time `json:"time,omitempty" bson:"time,omitempty"`
	MatchName      string     `json:"matchname,omitempty" bson:"matchname,omitempty"`
	MatchType      []string   `json:"type,omitempty" bson:"type,omitempty"`
	ModeName       string     `json:"mode_name,omitempty" bson:"mode_name,omitempty"`
	ModeDesc       string     `json:"mode_desc,omitempty" bson:"mode_desc,omitempty"`
	HandicapAmount string     `json:"handicap,omitempty" bson:"handicap,omitempty"`
	HandicapTeam   string     `json:"handicap_team,omitempty" bson:"handicap_team,omitempty"`
	RatioA         float64    `json:"ratioa,omitempty" bson:"ratioa"`
	RatioB         float64    `json:"ratiob,omitempty" bson:"ratiob"`
	Winner         string     `json:"winner,omitempty" bson:"winner,omitempty"`
	Status         string     `json:"status,omitempty" bson:"status,omitempty"`
	ScoreA         float64    `json:"scorea,omitempty" bson:"scorea"`
	ScoreB         float64    `json:"scoreb,omitempty" bson:"scoreb"`
	Note           string     `json:"note,omitempty" bson:"note,omitempty"`
	SeriesID       string     `json:"series_id,omitempty" bson:"series_id,omitempty"`
}

type F10kHistory struct {
	Name   string     `json:"name"`
	Kill   float64    `json:"kill,omitempty"`
	Death  float64    `json:"death,omitempty"`
	Winner string     `json:"winner,omitempty"`
	Time   *time.Time `json:"time,omitempty"`
}

type F10kResult struct {
	Name         string        `json:"name"`
	AverageKill  float64       `json:"avgkill"`
	AverageDeath float64       `json:"avgdeath"`
	RatioKill    float64       `json:"name"`
	TotalKill    float64       `json:"totalkill"`
	TotalDeath   float64       `json:"totaldeath"`
	Winrate      float64       `json:"winrate"`
	AverageOdds  float64       `json:"avgodds"`
	F10kHistory  []F10kHistory `json:"f10kHistory"`
}
