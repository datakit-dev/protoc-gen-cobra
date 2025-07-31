PROTO_SOURCES = $(wildcard */*.proto)
GO_SOURCES = $(patsubst %.proto,%.pb.go,$(PROTO_SOURCES))
GRPC_SOURCES = $(patsubst %.proto,%_grpc.pb.go,$(PROTO_SOURCES))
COBRA_SOURCES = $(patsubst %.proto,%.cobra.pb.go,$(PROTO_SOURCES))

all:
	protoc \
		--go_out=pb \
		--go-grpc_out=pb \
		--cobra_out=pb \
		$(wildcard pb/*.proto)
	protoc \
		-I. \
		--go_out=../../../.. \
		--go-grpc_out=../../../.. \
		--cobra_out=../../../.. \
		oneof/common/common.proto
	protoc \
		-Ioneof \
		--go_out=oneof \
		--go-grpc_out=oneof \
		--cobra_out=oneof \
		oneof/oneof.proto

clean:
	rm -f $(GO_SOURCES) $(GRPC_SOURCES) $(COBRA_SOURCES)

deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/datakit-dev/protoc-gen-cobra@latest
