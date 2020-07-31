module github.com/dictav/run-to-run/grpc/world

go 1.14

replace github.com/dictav/run-to-run/grpc/proto => ../proto

require (
	github.com/dictav/run-to-run/grpc/proto v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.31.0
)
