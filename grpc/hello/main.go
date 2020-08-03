package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	_log "log"
	"net"
	"os"
	"strings"

	"github.com/dictav/run-to-run/grpc/proto"

	cmeta "cloud.google.com/go/compute/metadata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var log = _log.New(os.Stderr, "", 0) //nolint:gochecknoglobals

//nolint:gomnd
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	var (
		err error

		ep  = ":" + port
		srv = grpc.NewServer()
	)

	s := &server{
		worldAddr: os.Getenv("RUN_WORLD_ADDR"),
		authToken: os.Getenv("AUTH_TOKEN"),
	}

	if s.authToken == "" {
		s.authToken, err = cmeta.Get(s.tokenURL())
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}
	}

	proto.RegisterHelloServer(srv, s)

	listen, err := net.Listen("tcp", ep)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	log.Println("grpc-hello is listener on ", ep)

	if err := srv.Serve(listen); err != nil {
		log.Println(err)
		os.Exit(2)
	}

	os.Exit(0)
}

type server struct {
	worldAddr string
	authToken string
}

func (s *server) tokenURL() string {
	cmps := strings.Split(s.worldAddr, ":")
	host := cmps[0]

	return "/instance/service-accounts/default/identity?audience=https://" + host + "/"
}

func (s *server) Hello(ctx context.Context, r *proto.HelloRequest) (*proto.HelloResponse, error) {
	con, err := dial(s.worldAddr, false)
	if err != nil {
		return nil, err
	}

	c := proto.NewWorldClient(con)
	req := &proto.WorldRequest{Text: "hello"}

	if s.authToken != "" {
		ctx = metadata.AppendToOutgoingContext(
			ctx,
			"authorization", "Bearer "+s.authToken,
		)
	}

	res, err := c.World(ctx, req)
	if err != nil {
		return nil, err
	}

	return &proto.HelloResponse{
		Text: res.Text,
	}, nil
}

func dial(addr string, insecure bool) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{}

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
