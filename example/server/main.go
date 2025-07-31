package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/datakit-dev/protoc-gen-cobra/example/pb"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer()
	pb.RegisterBankServer(srv, NewBank())
	pb.RegisterCacheServer(srv, NewCache())
	pb.RegisterCRUDServer(srv, NewCRUD())
	pb.RegisterCyclicalServer(srv, NewCyclical())
	pb.RegisterDeprecatedServer(srv, NewDeprecated())
	pb.RegisterNestedServer(srv, NewNested())
	pb.RegisterOneofServer(srv, NewOneof())
	pb.RegisterProto2Server(srv, NewProto2())
	pb.RegisterTimerServer(srv, NewTimer())
	pb.RegisterTypesServer(srv, NewTypes())
	err = srv.Serve(ln)
	if err != nil {
		log.Fatal(err)
	}
}
