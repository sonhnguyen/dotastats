package main

import (
	"dotastats"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/kardianos/osext"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

type dotastatsConfig struct {
	Port          string
	URI           string
	Dbname        string
	Collection    string
	IsDevelopment string
}

// App in main app
type App struct {
	router  *Router
	gp      globalPresenter
	logr    appLogger
	mongodb dotastats.Mongodb
	config  dotastatsConfig
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
			IsDevelopment: viper.GetString("isDevelopment"),
			Port:          viper.GetString("port"),
			URI:           viper.GetString("uri"),
			Dbname:        viper.GetString("dbname"),
			Collection:    viper.GetString("collection"),
		}
	} else {
		config = dotastatsConfig{
			IsDevelopment: os.Getenv("isDevelopment"),
			Port:          os.Getenv("PORT"),
			URI:           os.Getenv("uri"),
			Dbname:        os.Getenv("dbname"),
			Collection:    os.Getenv("collection"),
		}
	}

	mongo := dotastats.Mongodb{URI: config.URI, Dbname: config.Dbname, Collection: config.Collection}

	gp := globalPresenter{
		SiteName:    "dotastats",
		Description: "Api for native app",
		SiteURL:     "api.floatingcube.com",
	}

	return &App{
		router:  r,
		gp:      gp,
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
		fmt.Println("panicking")
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	r := NewRouter()
	logr := newLogger()
	a := SetupApp(r, logr, "")
	err = a.RunCrawlerAndSave()
	if err != nil {
		fmt.Errorf("error running crawler %s", err)
	}

	// common := alice.New(context.ClearHandler, a.loggingHandler, a.recoverHandler)
	// r.Get("/video/link", common.Then(a.Wrap(a.GetVideoByLinkHandler())))
	// r.Get("/video/id/:id", common.Then(a.Wrap(a.GetVideoByIdHandler())))
	// r.Get("/video/id/:id/subtitle", common.Then(a.Wrap(a.GetSubtitleByIDHandler())))
	// r.Get("/video", common.Then(a.Wrap(a.GetAllVideoHandler())))
	// r.Get("/video/random", common.Then(a.Wrap(a.GetRandomVideoHandler())))
	// r.Post("/video/:id", common.Then(a.Wrap(a.PostCommentByIdHandler())))

	// Add CORS support (Cross Origin Resource Sharing)
	handler := cors.Default().Handler(r)
	err = http.ListenAndServe(":"+a.config.Port, handler)
	if err != nil {
		fmt.Errorf("error on serve server %s", err)
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
