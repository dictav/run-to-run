package main

import (
	"context"
	"fmt"
	_log "log"
	"net"
	"os"

	"github.com/dictav/run-to-run/grpc/proto"

	"google.golang.org/grpc"
)

var log = _log.New(os.Stderr, "", 0)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ep := ":" + port
	srv := grpc.NewServer()

	s := &server{}

	proto.RegisterWorldServer(srv, s)

	listen, err := net.Listen("tcp", ep)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	log.Println("grpc-world is listensing onr", ep)

	if err := srv.Serve(listen); err != nil {
		log.Println(err)
		os.Exit(2)
	}

	os.Exit(0)
}

type server struct{}

func (s *server) World(ctx context.Context, r *proto.WorldRequest) (*proto.WorldResponse, error) {
	return &proto.WorldResponse{
		Text: fmt.Sprintf("%s world", r.Text),
	}, nil
}
