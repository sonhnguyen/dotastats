package dotastats

import (
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

func trimLine(line string) string {
	s := strings.TrimSpace(line)
	return s
}

func ratioProcess(ratio string) float64 {
	re := regexp.MustCompile("\\d+\\.\\d+")
	number := re.FindAllString(ratio, -1)
	if number != nil {
		result, err := strconv.ParseFloat(number[0], 64)
		if err != nil {
			return 0
		}
		return result
	}
	return 0
}
func winnerProcess(winner string) string {
	if len(winner) < 10 {
		return "TBD"
	}
	return winner[9:]
}

func processMatchesDota2BY(listMatches []string) ([]Match, error) {
	var result []Match
	for _, link := range listMatches {
		doc, err := goquery.NewDocument(link)
		if err != nil {
			return result, err
		}
		teamA := doc.Find("div.main-vs > div.op1 > div.title-opt > h3 > a:nth-child(2) > span").Text()
		teamB := doc.Find("div.main-vs > div.op2 > div.title-opt > h3 > a:nth-child(2) > span").Text()
		matchName := teamA + " vs " + teamB
		timeString := doc.Find("div.time-right").Text()

		layOut := "02 Jan 2006 at 15:04 MST"
		timeStamp, err := time.Parse(layOut, timeString)
		if err != nil {
			fmt.Errorf("error parsing date %s", err)
		}

		tournament := doc.Find("div.title2 > span.tt-right").Text()
		matchID := doc.Find("div.title2 > span.tt-left").Text()
		matchIDInt := matchIDProcess(matchID)
		bestOf := trimLine(doc.Find("div.kind-match").Text())
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
			ratioA := ratioProcess(s.Find("div.reward > div.left-reward > div.appid_570").Text())
			ratioB := ratioProcess(s.Find("div.reward > div.right-reward > div.appid_570").Text())
			matchType := []string{}
			matchType = append(matchType, trimLine(matchTypeList[i]))
			if matchType[0] == "Handicap" {
				matchType = append(matchType, trimLine(s.Find("div.full > div").Text()))
			}
			matchType = append(matchType, trimLine(doc.Find("div.txt-notice").Text()))

			winner := winnerProcess(s.Find("div.winner").Text())

			match := Match{MatchName: matchName, TeamA: teamA, TeamB: teamB, URL: link, Time: timeStamp, Tournament: tournament, MatchType: matchType, RatioA: ratioA, RatioB: ratioB, MatchID: matchIDInt, BestOf: bestOf, ScoreA: scoreA, ScoreB: scoreB, Winner: winner}
			result = append(result, match)
		})
	}
	return result, nil
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
	result, err = processMatchesDota2BY(listMatches)
	fmt.Println("done crawled %v", len(result))
	return result, err
}

func RunOldMatchesDota2BestYolo(startId int, endId int) ([]Match, error) {
	var result []Match
	var listMatches []string
	ROOT_URL := "http://dota2bestyolo.com/"
	for i := startId; i <= endId; i++ {
		listMatches = append(listMatches, ROOT_URL+"match/"+strconv.Itoa(i))
	}

	result, err := processMatchesDota2BY(listMatches)
	return result, err
}
