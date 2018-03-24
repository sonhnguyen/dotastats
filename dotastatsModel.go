package dotastats

import (
	"strings"
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

type Series struct {
	SeriesID string  `json:"series_id,omitempty"`
	Matches  []Match `json:"matches,omitempty"`
}

type PlayerInfo struct {
	FullName  string   `json:"fullname,omitempty" bson:"fullname,omitempty"`
	GameName  string   `json:"ingame_name,omitempty" bson:"ingame_name,omitempty"`
	Biography string   `json:"biography,omitempty" bson:"biography,omitempty"`
	Detail    string   `json:"detail,omitempty" bson:"detail,omitempty"`
	Links     []string `json:"links,omitempty" bson:"links,omitempty"`
	URL       string   `json:"url,omitempty" bson:"url,omitempty"`
}

func (p *PlayerInfo) FindTwitterID() string {
	for _, link := range p.Links {
		if i := strings.Index(link, "http://twitter.com/"); i != -1 {
			return link[i+19:]
		}
	}

	return ""
}

type TeamInfo struct {
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name,omitempty" bson:"name,omitempty"`
	NameSlug string        `json:"slug,omitempty" bson:"slug,omitempty"`
	Game     string        `json:"game,omitempty" bson:"game,omitempty"`
	Region   string        `json:"region,omitempty" bson:"region,omitempty"`
	Players  []PlayerInfo  `json:"players,omitempty" bson:"players,omitempty"`
	Overview string        `json:"overview,omitempty" bson:"overview,omitempty"`
	History  string        `json:"history,omitempty" bson:"history,omitempty"`
	Logo     string        `json:"logo,omitempty" bson:"logo,omitempty"`
	URL      string        `json:"url,omitempty" bson:"url,omitempty"`
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
	RatioA         float64    `json:"ratioa" bson:"ratioa"`
	RatioB         float64    `json:"ratiob" bson:"ratiob"`
	Winner         string     `json:"winner,omitempty" bson:"winner,omitempty"`
	Status         string     `json:"status,omitempty" bson:"status,omitempty"`
	ScoreA         float64    `json:"scorea" bson:"scorea"`
	ScoreB         float64    `json:"scoreb" bson:"scoreb"`
	Note           string     `json:"note,omitempty" bson:"note,omitempty"`
	SeriesID       string     `json:"series_id,omitempty" bson:"series_id,omitempty"`
}

type PicksBans struct {
	IsPick  bool `json:"is_pick,omitempty" bson:"is_pick,omitempty"`
	HeroID  int  `json:"hero_id,omitempty" bson:"hero_id,omitempty"`
	Team    int  `json:"team,omitempty" bson:"team,omitempty"`
	Order   int  `json:"ord,omitempty" bson:"ord,omitempty"`
	MatchID int  `json:"match_id,omitempty" bson:"match_id,omitempty"`
}

type OpenDotaMatch struct {
	Id             bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	MatchID        int           `json:"match_id" bson:"match_id"`
	Duration       int           `json:"duration" bson:"duration"`
	StartTime      *time.Time    `json:"start_time" bson:"start_time"`
	RadiantTeamID  int           `json:"radiant_team_id" bson:"radiant_team_id"`
	RadiantName    string        `json:"radiant_name" bson:"radiant_name"`
	RadiantTag     string        `json:"radiant_tag" bson:"radiant_tag"`
	RadiantLogoURL string        `json:"radiant_logo_url" bson:"radiant_logo_url"`
	DireTeamID     int           `json:"dire_team_id" bson:"dire_team_id"`
	DireName       string        `json:"dire_name" bson:"dire_name"`
	DireTag        string        `json:"dire_tag" bson:"dire_tag"`
	DireLogoURL    string        `json:"dire_logo_url" bson:"logo_url"`
	LeagueID       int           `json:"leagueid" bson:"leagueid"`
	LeagueName     string        `json:"league_name" bson:"league_name"`
	SeriesID       int           `json:"series_id" bson:"series_id"`
	SeriesType     int           `json:"series_type" bson:"series_type"`
	RadiantScore   int           `json:"radiant_score" bson:"radiant_score"`
	DireScore      int           `json:"dire_score" bson:"dire_score"`
	RadiantWin     bool          `json:"radiant_win" bson:"radiant_win"`
	PicksBans      []PicksBans   `json:"picks_bans" bson:"picks_bans"`
}

type F10kResult struct {
	Name         string  `json:"name"`
	AverageKill  float64 `json:"avgkill"`
	AverageDeath float64 `json:"avgdeath"`
	RatioKill    float64 `json:"ratiokill"`
	TotalKill    float64 `json:"totalkill"`
	TotalDeath   float64 `json:"totaldeath"`
	Winrate      float64 `json:"winrate"`
	AverageOdds  float64 `json:"avgodds"`
	Matches      []Match `json:"matches"`
}

type FBResult struct {
	Name    string  `json:"name"`
	Winrate float64 `json:"winrate"`
	Matches []Match `json:"matches"`
}

type TwitterCreateListRequest struct {
	Name        string `json:"name"`
	Mode        string `json:"mode"`
	Description string `json:"description"`
}

type TwitterAddToListRequest struct {
	OwnerScreenName string `json:"owner_screen_name"`
	Slug            string `json:"slug"`
	ScreenName      string `json:"screen_name"`
}

type TwitterRemoveListRequest struct {
	OwnerScreenName string `json:"owner_screen_name"`
	Slug            string `json:"slug"`
}

type TwitterGetListResponse struct {
	Lists []TwitterList `json:"lists"`
}

type TwitterList struct {
	Slug string `json:"slug"`
}

type Feedback struct {
	Name     string    `json:"name"`
	Feedback string    `json:"feedback"`
	Time     time.Time `json:"time"`
}

type User struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	RegisterKey string `json:"register_key"`
}

type Session struct {
	Email      string `json:"email"`
	SessionKey string `json:"session_key"`
}
