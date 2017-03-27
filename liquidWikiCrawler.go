package dotastats

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	slugify "github.com/metal3d/go-slugify"
)

const (
	LiquidBaseURL  = "http://wiki.teamliquid.net"
	DotaLiquidTeam = "http://wiki.teamliquid.net/dota2/Portal:Teams"
	CSGOLiquidTeam = "http://wiki.teamliquid.net/counterstrike/Portal:Teams"
)

func runCrawlerLiquidDota() ([]TeamInfo, error) {
	var result []TeamInfo
	teamList, err := goquery.NewDocument(DotaLiquidTeam)
	if err != nil {
		return []TeamInfo{}, err
	}

	teamList.Find("#mw-content-text > div:nth-child(5) > div.template-box").Each(func(i int, s *goquery.Selection) {
		s.Find("h3").Each(func(i int, regionSelect *goquery.Selection) {
			region := regionSelect.Find("span.mw-headline").Text()
			regionTeam := TeamInfo{Region: region}

			regionSelect.Next().Find("li").Each(func(i int, teamSelect *goquery.Selection) {
				team := regionTeam
				// For each item found, get the band and title
				url, exists := teamSelect.Find("span.team-template-text a").Attr("href")
				if exists {
					team.URL = LiquidBaseURL + url
					team.Name = teamSelect.Find("span.team-template-text a").Text()
					team.NameSlug = slugify.Marshal(team.Name)
					team.Game = "dota"
					result = append(result, team)
				}
			})
		})
	})
	for index, _ := range result {
		value := &result[index]

		teamDetail, err := goquery.NewDocument(value.URL)
		if err != nil {
			continue
		}
		logo, ok := teamDetail.Find(".infobox-image img").Attr("src")
		if ok {
			value.Logo = logo
		}
		teamDetail.Find(".activesquad .table.table-striped").First().Find("tbody > tr").Each(func(i int, playerSelect *goquery.Selection) {
			var player PlayerInfo
			player.GameName = playerSelect.Find("td:nth-child(2) a").Text()
			player.FullName = playerSelect.Find("td:nth-child(3)").Text()
			playerURL, ok := playerSelect.Find("td:nth-child(2) a").Attr("href")
			if ok && !playerSelect.Find("td:nth-child(2) a").HasClass("new") {
				player.URL = LiquidBaseURL + playerURL
				playerDoc, err := goquery.NewDocument(player.URL)
				if err != nil {
					return
				}
				playerDoc.Find(".infobox-center.infobox-icons a").Each(func(i int, playerLinks *goquery.Selection) {
					link, ok := playerLinks.Attr("href")
					if ok {
						player.Links = append(player.Links, link)

					}
				})
			}
			value.Players = append(value.Players, player)
		})
	}
	return result, nil
}
func runCrawlerLiquidCSGO() ([]TeamInfo, error) {
	var result []TeamInfo
	teamList, err := goquery.NewDocument(CSGOLiquidTeam)
	if err != nil {
		return []TeamInfo{}, err
	}

	teamList.Find("#mw-content-text > div > div > div > table:nth-child(5) > tbody > tr > td").Each(func(i int, s *goquery.Selection) {
		s.Find("h3").Each(func(i int, regionSelect *goquery.Selection) {
			region := regionSelect.Find("span.mw-headline").Text()
			regionTeam := TeamInfo{Region: region}

			regionSelect.Next().Find("li").Each(func(i int, teamSelect *goquery.Selection) {
				team := regionTeam
				// For each item found, get the band and title
				url, exists := teamSelect.Find("span.team-template-text a").Attr("href")
				if exists {
					team.URL = LiquidBaseURL + url
					team.Name = teamSelect.Find("span.team-template-text a").Text()
					team.NameSlug = slugify.Marshal(team.Name)
					team.Game = "csgo"
					result = append(result, team)
				}
			})
		})
	})
	for index, _ := range result {
		value := &result[index]
		teamDetail, err := goquery.NewDocument(value.URL)
		if err != nil {
			continue
		}
		logo, ok := teamDetail.Find(".infobox-image img").Attr("src")
		if ok {
			value.Logo = logo

		}
		teamDetail.Find(".table.table-striped").First().Find("tbody > tr").Each(func(i int, playerSelect *goquery.Selection) {
			var player PlayerInfo
			player.GameName = playerSelect.Find("td:nth-child(2) a").Text()
			player.FullName = playerSelect.Find("td:nth-child(3)").Text()
			playerURL, ok := playerSelect.Find("td:nth-child(2) a").Attr("href")
			if ok && !playerSelect.Find("td:nth-child(2) a").HasClass("new") {
				player.URL = LiquidBaseURL + playerURL
				playerDoc, err := goquery.NewDocument(player.URL)
				if err != nil {
					return
				}
				playerDoc.Find(".infobox-center.infobox-icons a").Each(func(i int, playerLinks *goquery.Selection) {
					link, ok := playerLinks.Attr("href")
					if ok {
						player.Links = append(player.Links, link)
					}
				})
			}
			value.Players = append(value.Players, player)
		})
	}
	return result, nil
}

func RunCrawlerLiquidTeam(game string) ([]TeamInfo, error) {
	var result []TeamInfo
	var err error
	if game == "dota" {
		result, err = runCrawlerLiquidDota()
		fmt.Println("crawled team info dota", len(result))
		if err != nil {
			return []TeamInfo{}, err
		}
	} else if game == "csgo" {
		result, err = runCrawlerLiquidCSGO()
		fmt.Println("crawled team info csgo", len(result))
		if err != nil {
			return []TeamInfo{}, err
		}
	}
	return result, nil

}
