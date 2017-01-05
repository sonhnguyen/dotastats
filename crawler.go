package dotastats

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func scoreProcess(score string) []int {
	var result []int
	score = strings.TrimSpace(score)
	scoreArr := strings.Split(score, ":")
	for _, scoreString := range scoreArr {
		if scoreString != " " {
			scoreInt, err := strconv.Atoi(strings.TrimSpace(scoreString))
			if err != nil {
				return []int{}
			}
			result = append(result, scoreInt)
		}
	}
	return result
}

func matchIDProcess(matchID string) int {
	re := regexp.MustCompile("[^0-9]")
	result, err := strconv.Atoi(re.ReplaceAllString(matchID, ""))
	if err != nil {
		return 0
	}
	return result

}
func RunCrawlerDota2BestYolo() ([]Match, error) {
	var result []Match
	var listMatches []string
	ROOT_URL := "http://dota2bestyolo.com"
	doc, err := goquery.NewDocument(ROOT_URL)
	if err != nil {
		return []Match{}, err
	}

	doc.Find("div.blk2 div.view-btn > a").Each(func(i int, s *goquery.Selection) {
		matchLink, _ := s.Attr("href")
		if matchLink != "" {
			listMatches = append(listMatches, ROOT_URL+matchLink)
		}
	})

	for _, link := range listMatches {
		doc, err := goquery.NewDocument(link)
		if err != nil {
			return result, err
		}
		teamA := doc.Find("div.main-vs > div.op1 > div.title-opt > h3 > a:nth-child(2) > span").Text()
		teamB := doc.Find("div.main-vs > div.op2 > div.title-opt > h3 > a:nth-child(2) > span").Text()
		timeString := doc.Find("time-right").Text()

		layOut := "02 Jan 2006 at 15:04 MST"
		timeStamp, err := time.Parse(layOut, timeString)
		if err != nil {
			fmt.Errorf("error parsing date %s", err)
		}

		tournament := doc.Find("div.title2 > span.tt-right").Text()
		matchID := doc.Find("div.title2 > span.tt-left").Text()
		matchIDInt := matchIDProcess(matchID)
		bestOf := doc.Find("div.kind-match").Text()
		score := doc.Find("div.main-vs div.vs span").Text()

		scoreArray := scoreProcess(score)
		var scoreA int
		var scoreB int
		if len(scoreArray) > 0 {
			scoreA = scoreArray[0]
			scoreB = scoreArray[1]
		}

		matchTypeList := []string{}
		doc.Find("li[role='presentation']").Each(func(i int, s *goquery.Selection) {
			matchTypeList = append(matchTypeList, s.Find("a").Text())
		})

		doc.Find("div.match-bk1 > div").Each(func(i int, s *goquery.Selection) {
			ratioA := s.Find("div.reward > div.left-reward > div.appid_570").Text()
			ratioB := s.Find("div.reward > div.right-reward > div.appid_570").Text()
			matchType := []string{}
			matchType = append(matchType, matchTypeList[i])
			if matchType[0] == "Handicap" {
				matchType = append(matchType, s.Find("div.full > div").Text())
			}
			matchType = append(matchType, doc.Find("div.txt-notice").Text())

			winner := s.Find("div.winner").Text()

			match := Match{TeamA: teamA, TeamB: teamB, URL: link, Time: timeStamp, Tournament: tournament, MatchType: matchType, RatioA: ratioA, RatioB: ratioB, MatchID: matchIDInt, BestOf: bestOf, ScoreA: scoreA, ScoreB: scoreB, Winner: winner}
			result = append(result, match)
		})
	}
	b, _ := json.Marshal(result)
	fmt.Println("%v", string(b[:]))
	return result, nil
}
