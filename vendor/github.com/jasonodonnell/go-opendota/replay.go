package opendota

import (
	"net/http"

	"github.com/dghubble/sling"
)

func newReplayService(sling *sling.Sling) *ReplayService {
	return &ReplayService{
		sling: sling.Path("replays"),
	}
}

// ReplayService provides a method for accesing replay data.
type ReplayService struct {
	sling *sling.Sling
}

type replayParam struct {
	MatchID []int `url:"match_id"`
}

// Replay represents a Dota 2 replay.
type Replay struct {
	MatchID    int64 `json:"match_id"`
	Cluster    int   `json:"cluster"`
	ReplaySalt int   `json:"replay_salt"`
	SeriesID   int   `json:"series_id"`
	SeriesType int   `json:"series_type"`
}

// Replays takes an array of  Match IDs and returns replays for those
// matches.
// https://docs.opendota.com/#tag/replays%2Fpaths%2F~1replays%2Fget
func (s *ReplayService) Replays(matchID []int) ([]Replay, *http.Response, error) {
	params := &replayParam{}
	params.MatchID = matchID
	replays := new([]Replay)
	apiError := new(APIError)
	resp, err := s.sling.New().QueryStruct(params).Receive(replays, apiError)
	return *replays, resp, relevantError(err, *apiError)
}
