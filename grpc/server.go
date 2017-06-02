package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/axiomzen/zenauth/config"
	"github.com/axiomzen/zenauth/data"
	"github.com/axiomzen/zenauth/protobuf"

	google_grpc "google.golang.org/grpc"
)

type Server struct {
	Config *config.ZENAUTHConfig
	DAL    data.ZENAUTHProvider
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", s.Config.GRPCPort))
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := google_grpc.NewServer()
	protobuf.RegisterAuthServer(grpcServer, &Auth{
		Config: s.Config,
		DAL:    s.DAL,
	})
	log.Printf("Starting GRPC Server on Port %v", s.Config.GRPCPort)
	return grpcServer.Serve(ln)
}
