package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"

	"github.com/dictav/run-to-run/grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ep := ":" + port
	srv := grpc.NewServer()

	s := &server{worldAddr: os.Getenv("RUN_WORLD_ADDR")}

	proto.RegisterHelloServer(srv, s)

	listen, err := net.Listen("tcp", ep)
	if err != nil {
		println(err.Error())
		os.Exit(2)
	}

	println("Starting: gRPC Listener", ep)

	if err := srv.Serve(listen); err != nil {
		println(err.Error())
		os.Exit(2)
	}

	os.Exit(0)
}

type server struct {
	worldAddr string
}

func (s *server) Hello(ctx context.Context, r *proto.HelloRequest) (*proto.HelloResponse, error) {
	con, err := dial(s.worldAddr, false)
	if err != nil {
		return nil, err
	}

	c := proto.NewWorldClient(con)
	req := &proto.WorldRequest{Text: "hello"}

	res, err := c.World(ctx, req)
	if err != nil {
		return nil, err
	}

	return &proto.HelloResponse{
		Text: res.Text,
	}, nil
}

func dial(addr string, insecure bool) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	}

	if insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}

		creds := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})

		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	return grpc.Dial(addr, opts...)
}