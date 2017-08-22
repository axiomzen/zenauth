package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/axiomzen/zenauth/config"
	"github.com/axiomzen/zenauth/data"
	"github.com/axiomzen/zenauth/helpers"
	"github.com/axiomzen/zenauth/protobuf"
	newrelic "github.com/newrelic/go-agent"

	google_grpc "google.golang.org/grpc"
)

type Server struct {
	Config      *config.ZENAUTHConfig
	DAL         data.ZENAUTHProvider
	Log         *logrus.Entry
	NewRelicApp newrelic.Application
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", s.Config.GRPCPort))
	if err != nil {
		log.Fatal(err)
	}
	if s.Config.NewRelicEnabled {
		cfg := newrelic.NewConfig(s.Config.NewRelicName, s.Config.NewRelicKey)
		nrApp, err := newrelic.NewApplication(cfg)
		if err != nil {
			log.Fatal(err)
		} else {
			s.NewRelicApp = nrApp
		}
	}
	authServer := &Auth{
		Config: s.Config,
		DAL:    s.DAL,
		Log:    s.Log.WithField("GRPC Service", "Auth"),
		monitor: helpers.Monitor{
			Enabled:     s.Config.NewRelicEnabled,
			NewRelicApp: s.NewRelicApp,
		},
	}
	grpcServer := google_grpc.NewServer(google_grpc.UnaryInterceptor(authServer.AuthUnaryInterceptor))
	protobuf.RegisterAuthServer(grpcServer, authServer)
	log.Printf("Starting GRPC Server on Port %v", s.Config.GRPCPort)
	return grpcServer.Serve(ln)
}
