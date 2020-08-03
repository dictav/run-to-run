package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"os"

	"github.com/dictav/run-to-run/grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var (
	addr     = flag.String("addr", "localhost:8080", "host")
	insecure = flag.Bool("insecure", false, "use insecure request")
)

func main() {
	flag.Parse()

	opts, err := getDialOpts(*insecure)
	if err != nil {
		println(err.Error())
		os.Exit(2)
	}

	println("Connecting to gRPC Service", *addr)

	conn, err := grpc.Dial(*addr, opts...)
	if err != nil {
		println(err.Error())
		os.Exit(2)
	}

	defer conn.Close()

	client := proto.NewHelloClient(conn)
	ctx := metadata.AppendToOutgoingContext(
		context.Background(),
		"authorization", "Bearer "+os.Getenv("AUTH_TOKEN"),
	)
	req := proto.HelloRequest{}

	res, err := client.Hello(ctx, &req)
	if err != nil {
		println(err.Error())
		os.Exit(2)
	}

	println(res.Text)
}

func getDialOpts(insecure bool) ([]grpc.DialOption, error) {
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

	return opts, nil
}
