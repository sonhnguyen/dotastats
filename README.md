A dota app that crawl dota 2 matches result and get informative stuff

A match JSON:
```JSON
{
    "_id" : ObjectId("588a5498ca9cb54728cad292"),
    "url" : "http://www.vpgame.com/match/10113208",
    "teama_id" : "100000138",
    "teamb_id" : "26160",
    "teama" : "Vici Gaming",
    "teamb" : "WAY",
    "teama_short" : "VG",
    "teamb_short" : "WAY",
    "tournament" : "Dota 2 Asia Championships",
    "game" : "dota",
    "bestof" : "BO1",
    "matchid" : "10113208",
    "time" : ISODate("2017-01-14T05:47:00.000Z"),
    "matchname" : "Vici Gaming vs WAY, Game1 10kills",
    "mode_name" : "Game1 10kills",
    "handicap" : "0",
    "ratioa" : 0.33,
    "ratiob" : 2.74,
    "winner" : "Vici Gaming",
    "status" : "Settled",
    "scorea" : 10.0,
    "scoreb" : 4.0
}
```

A F10k report:
```JSON
{
	"avgkill":8.333333333333334,
	"avgdeath":7.5,
	"totalkill":100,
	"totaldeath":90,
	"winrate":0.5,
	"avgodds":0.8483333333333333,
	"enemy":["og","og","secret","secret","secret","faceless","faceless","faceless","newbee","newbee","wings","wings"]
}
```

- GET /team/:name?limit=10 (default no limit): Return matches of a team, sorted by most recent
    - URL params:
        - Limit (default 100)
        - Skip (default 0)
        - Fields (default show all fields), Specify fields to returns, format as array of string (&fields=teama,teama_id,ratioa,...)

- GET /team/:name/f10k?limit=10 (default no limit): Return f10k matches of a team, sorted by most recent
    - URL params:
        - Limit (default 100)
        - Skip (default 0)
        - Fields (default show all fields), Specify fields to returns, format as array of string (&fields=teama,teama_id,ratioa,...)
        
- GET /f10k/:name?limit=10 (default no limit): Return F10k report of a team.
    - URL params:
        - Limit (default 100)
        - Skip (default 0)

- GET /match
    - URL params:
        - Limit (default 100)
        - Skip (default 0)
        - Status (default all), could be `open`, `closed` or `live`
        - Fields (default show all fields), Specify fields to returns, format as array of string (&fields=teama,teama_id,ratioa,...)
