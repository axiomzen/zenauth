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
	"github.com/axiomzen/zenauth/grpc"
	pg "gopkg.in/pg.v4"
	"gopkg.in/tylerb/graceful.v1"
)

func main() {
	log.Infof(os.Getenv("ZENAUTH_ENVIRONMENT"))

	// set just in case someone has go 1.4
	runtime.GOMAXPROCS(runtime.NumCPU())

	nullformat.SetTimeFormat(constants.TimeFormat)

	log.SetFormatter(&log.TextFormatter{})

	// apparently this is the default now (we really should fork this repo)
	//uuid.SwitchFormat(uuid.CleanHyphen)

	conf, err := config.Get()

	if err != nil {
		// die, we are not configured properly
		log.Fatal(err.Error())
	}

	switch conf.Environment {
	case constants.EnvironmentStaging, constants.EnvironmentProduction, constants.EnvironmentDevelopment:
		data.Migrate(conf)
		log.SetFormatter(&log.JSONFormatter{})
		fallthrough
	case constants.EnvironmentTest:
		// set seed
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

	// Error channel for multiple servers
	errChn := make(chan error)

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
	go func() {
		errChn <- srv.ListenAndServe()
	}()

	// Runs the GRPC server
	grpcServer := grpc.Server{
		Config: conf,
		DAL:    dataP,
		Log:    log.WithField("server", "grpc"),
	}
	go func() {
		errChn <- grpcServer.ListenAndServe()
	}()

	log.Fatal(<-errChn)
}
