package main

import (
	"context"

	"github.com/datakit-dev/protoc-gen-cobra/example/pb"
)

type Cyclical struct {
	pb.UnimplementedCyclicalServer
}

func NewCyclical() *Cyclical {
	return &Cyclical{}
}

func (*Cyclical) Test(_ context.Context, req *pb.Foo) (*pb.Bar, error) {
	if req.Bar1 != nil {
		return req.Bar1, nil
	} else if req.Bar2 != nil {
		return req.Bar2, nil
	} else {
		return &pb.Bar{
			Foo1:  req,
			Foo2:  req,
			Value: "hello",
		}, nil
	}
}
