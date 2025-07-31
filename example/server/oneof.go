package main

import (
	"context"

	"github.com/datakit-dev/protoc-gen-cobra/example/pb"
)

type Oneof struct {
	pb.UnimplementedOneofServer
}

func NewOneof() *Oneof {
	return &Oneof{}
}

func (*Oneof) Fetch(_ context.Context, req *pb.FetchRequest) (*pb.FetchResponse, error) {
	return &pb.FetchResponse{Value: req.GetOption1() + req.GetOption2() + req.GetOption3()}, nil
}

func (*Oneof) FetchNested(_ context.Context, req *pb.FetchNestedRequest) (*pb.FetchResponse, error) {
	l2 := req.GetL0().GetL1().GetL2()
	return &pb.FetchResponse{Value: l2.GetOption1() + l2.GetOption2() + l2.GetOption3()}, nil
}
