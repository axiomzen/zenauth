package main

import (
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	logger "log"

	log "github.com/Sirupsen/logrus"
	nullformat "github.com/axiomzen/null/format"
	"github.com/axiomzen/zenauth/config"
	"github.com/axiomzen/zenauth/constants"
	"github.com/axiomzen/zenauth/data"
	"github.com/joho/godotenv"
	pg "gopkg.in/pg.v4"
	"gopkg.in/tylerb/graceful.v1"
)

func main() {
	// set local dev env
	if strings.EqualFold(os.Getenv("ZENAUTH_ENVIRONMENT"), constants.EnvironmentDevelopment) {
		if err := godotenv.Load(); err != nil {
			// Do not fail here, incase they've manually loaded env variables
			// Will fail out at config.Get() if any require env variables not set
			log.WithError(err).Warn("error loading godotenv")
		}
	}
	// set just in case someone has go 1.4
	runtime.GOMAXPROCS(runtime.NumCPU())

	nullformat.SetTimeFormat(constants.TimeFormat)

	log.SetFormatter(&log.JSONFormatter{})

	// apparently this is the default now (we really should fork this repo)
	//uuid.SwitchFormat(uuid.CleanHyphen)

	conf, err := config.Get()

	if err != nil {
		// die, we are not configured properly
		log.Fatal(err.Error())
	}

	// set seed
	switch conf.Environment {
	case constants.EnvironmentStaging, constants.EnvironmentTest, constants.EnvironmentProduction:
		rand.Seed(time.Now().UTC().UnixNano())
	default:
	}

	// set query logger
	if conf.LogQueries {
		pg.SetQueryLogger(logger.New(os.Stdout, "", logger.LstdFlags))
	}

	// set logging level
	switch strings.ToLower(conf.LogLevel) {
	default:
		fallthrough
	case log.DebugLevel.String():
		log.SetLevel(log.DebugLevel)
	case log.InfoLevel.String():
		log.SetLevel(log.InfoLevel)
	case log.WarnLevel.String():
		log.SetLevel(log.WarnLevel)
	case log.ErrorLevel.String():
		log.SetLevel(log.ErrorLevel)
	case log.FatalLevel.String():
		log.SetLevel(log.FatalLevel)
	case log.PanicLevel.String():
		log.SetLevel(log.PanicLevel)
	}

	// database
	dataP, dataErr := data.Get(conf)

	if dataErr != nil {
		// die, we can't connect to the database
		log.Fatal(dataErr.Error())
	}

	// make sure to close the database connection pool when we exit
	defer dataP.Close()

	router := InitRouter(conf)

	srv := &graceful.Server{
		// Time to allow for active requests to complete
		Timeout: conf.DrainAndDieTimeout,

		Server: &http.Server{
			Addr:         ":" + strconv.FormatInt(int64(conf.Port), 10),
			Handler:      router,
			ReadTimeout:  conf.TransportReadTimeout,
			WriteTimeout: conf.TransportWriteTimeout,
			//MaxHeaderBytes: 1 << 20,
		},
	}

	log.Println("Starting Server on Port " + strconv.FormatInt(int64(conf.Port), 10))
	log.Fatal(srv.ListenAndServe())
}
