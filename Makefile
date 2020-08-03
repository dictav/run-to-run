HTTP_HELLO_URL := https://http-hello-iljqu5cu6q-uc.a.run.app
GRPC_HELLO_ADDR := grpc-hello-iljqu5cu6q-uc.a.run.app:443
AUTH_TOKEN := $(shell gcloud auth print-identity-token)

grpchello: cmd/grpchello/main.go
	cd cmd/grpchello; env AUTH_TOKEN=$(AUTH_TOKEN) go run main.go -addr=$(GRPC_HELLO_ADDR)

httphello:
	curl -iH "Authorization: Bearer $(AUTH_TOKEN)" $(HTTP_HELLO_URL)
