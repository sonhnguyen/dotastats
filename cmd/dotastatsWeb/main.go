package main

import (
	"dotastats"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
	"github.com/kardianos/osext"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	cron "gopkg.in/robfig/cron.v2"
)

type dotastatsConfig struct {
	Port               string
	URI                string
	Dbname             string
	Collection         string
	CollectionTeam     string
	CollectionProMatch string
	CollectionFeedback string
	CollectionUser     string
	CollectionSession  string
	CookieSecretKey    string
	IsDevelopment      string
	RegisterKey        string
}

// App in main app
type App struct {
	router  *Router
	gp      globalPresenter
	logr    appLogger
	mongodb dotastats.Mongodb
	config  dotastatsConfig
	store   *sessions.CookieStore
}

// globalPresenter contains the fields neccessary for presenting in all templates
type globalPresenter struct {
	SiteName    string
	Description string
	SiteURL     string
}

// TODO localPresenter if we have using template
func SetupApp(r *Router, logger appLogger, templateDirectoryPath string) *App {
	var config dotastatsConfig
	if viper.GetBool("isDevelopment") {
		config = dotastatsConfig{
			IsDevelopment:      viper.GetString("isDevelopment"),
			Port:               viper.GetString("port"),
			URI:                viper.GetString("uri"),
			Dbname:             viper.GetString("dbname"),
			Collection:         viper.GetString("collection"),
			CollectionTeam:     viper.GetString("collection-team"),
			CollectionProMatch: viper.GetString("collection-pro-match"),
			CollectionFeedback: viper.GetString("collection-feedback"),
			CollectionUser:     viper.GetString("collection-user"),
			CollectionSession:  viper.GetString("collection-session"),
			CookieSecretKey:    viper.GetString("cookie-secret-key"),
			RegisterKey:        viper.GetString("register-key"),
		}
	} else {
		config = dotastatsConfig{
			IsDevelopment:      os.Getenv("isDevelopment"),
			Port:               os.Getenv("PORT"),
			URI:                os.Getenv("uri"),
			Dbname:             os.Getenv("dbname"),
			Collection:         os.Getenv("collection"),
			CollectionTeam:     os.Getenv("collection-team"),
			CollectionProMatch: os.Getenv("collection-pro-match"),
			CollectionFeedback: os.Getenv("collection-feedback"),
			CollectionUser:     os.Getenv("collection-user"),
			CollectionSession:  os.Getenv("collection-session"),
			CookieSecretKey:    os.Getenv("cookie-secret-key"),
			RegisterKey:        os.Getenv("register-key"),
		}
	}

	if viper.GetBool("isLocal") {
		config.URI = viper.GetString("uriLocal")
	}

	mongo := dotastats.Mongodb{
		URI:                config.URI,
		Dbname:             config.Dbname,
		Collection:         config.Collection,
		CollectionTeam:     config.CollectionTeam,
		CollectionProMatch: config.CollectionProMatch,
		CollectionFeedback: config.CollectionFeedback,
		CollectionUser:     config.CollectionUser,
		CollectionSession:  config.CollectionSession,
	}

	gp := globalPresenter{
		SiteName:    "dotastats",
		Description: "Api",
		SiteURL:     "wtf",
	}

	return &App{
		router:  r,
		gp:      gp,
		store:   sessions.NewCookieStore([]byte(config.CookieSecretKey)),
		logr:    logger,
		config:  config,
		mongodb: mongo,
	}
}

func main() {
	pwd, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatalf("cannot retrieve present working directory: %i", 0600, nil)
	}

	err = LoadConfiguration(pwd)
	if err != nil && os.Getenv("PORT") == "" {
		log.Panicln("panicking, Fatal error config file: %s", err)
	}

	r := NewRouter()
	logr := newLogger()
	a := SetupApp(r, logr, "")
	// Add CORS support (Cross Origin Resource Sharing)
	corsSetting := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://f10k.herokuapp.com", "http://dotastats.me", "http://www.dotastats.me"},
		AllowCredentials: true,
	})
	handler := corsSetting.Handler(r)
	if a.config.IsDevelopment == "true" {
		handler = cors.Default().Handler(r)
	}

	common := alice.New(context.ClearHandler, a.loggingHandler, a.recoverHandler)
	authenticate := common.Append(a.UserMiddlewareGenerator, a.authMiddleware)
	r.Get("/team-info/:slug", common.Then(a.Wrap(a.GetTeamInfoHandler())))
	r.Get("/team/:name", common.Then(a.Wrap(a.GetTeamMatchesHandler())))
	r.Get("/history", common.Then(a.Wrap(a.GetTeamHistoryHandler())))
	r.Get("/team/:name/f10k", common.Then(a.Wrap(a.GetTeamF10kMatchesHandler())))
	r.Get("/team/:name/fb", common.Then(a.Wrap(a.GetTeamFBMatchesHandler())))
	r.Get("/matches", common.Then(a.Wrap(a.GetMatchesListHandler())))
	r.Get("/matches/:id", common.Then(a.Wrap(a.GetMatchByIDHandler())))
	r.Get("/crawl", common.Then(a.Wrap(a.GetCustomCrawlHandler())))
	r.Get("/crawlTeamInfo", common.Then(a.Wrap(a.GetCrawlTeamInfoHandler())))
	r.Get("/create-twitter-list", common.Then(a.Wrap(a.CreateAllTwitterList())))
	r.Get("/remove-twitter-list", common.Then(a.Wrap(a.RemoveAllTwitterList())))
	r.Get("/feedback", authenticate.Then(a.Wrap(a.GetFeedback())))
	r.Post("/feedback", common.Then(a.Wrap(a.PostFeedback())))

	r.Post("/login", common.Then(a.Wrap(a.LoginPostHandler())))
	r.Post("/register", common.Then(a.Wrap(a.RegisterPostHandler())))

	c := cron.New()
	_, err = c.AddFunc("@every 1s", func() {
		err = a.RunCrawlerAndSave()
		if err != nil {
			log.Println("error running crawler %s", err)
		}
		err = a.RunCrawlerOpenDotaProMatchesAndSave()
		if err != nil {
			log.Println("error running crawler %s", err)
		}
	})
	if err != nil {
		log.Println("error on cron job %s", err)
	}
	_, err = c.AddFunc("@weekly", func() {
		a.RunCrawlerTeamInfoAndSave()
		if err != nil {
			log.Println("error running crawler %s", err)
		}
	})
	if err != nil {
		log.Println("error on cron job %s", err)
	}
	c.Start()
	err = http.ListenAndServe(":"+a.config.Port, handler)
	if err != nil {
		log.Println("error on serve server %s", err)
	}
}

func LoadConfiguration(pwd string) error {
	viper.SetConfigName("dotastats-config")
	viper.AddConfigPath(pwd)
	devPath := pwd[:len(pwd)-3] + "src/dotastats/cmd/dotastatsweb/"
	_, file, _, _ := runtime.Caller(1)
	configPath := path.Dir(file)
	viper.AddConfigPath(devPath)
	viper.AddConfigPath(configPath)
	return viper.ReadInConfig() // Find and read the config file
}
