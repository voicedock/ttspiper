package main

import (
	"fmt"
	grpcapi "github.com/test/test/internal/api/grpc"
	ttsv1 "github.com/test/test/internal/api/grpc/gen/test/tts/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	srv := grpcapi.NewServerTts()
	fmt.Println(srv)

	lis, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		fmt.Printf("failed to listen GRPC server: %s\n", err)
	}

	s := grpc.NewServer()
	ttsv1.RegisterTtsAPIServer(s, srv)
	reflection.Register(s)
	s.Serve(lis)
}
