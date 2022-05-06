package main

import (
	"github.com/gudn/lesslog/proto"
)

type lesslogServer struct {
	proto.UnimplementedLesslogServer
}

func Build() *lesslogServer {
	return &lesslogServer{}
}
