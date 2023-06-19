package main

import (
	"fmt"
	grpcapi "github.com/voicedock/ttspiper/internal/api/grpc"
	ttsv1 "github.com/voicedock/ttspiper/internal/api/grpc/gen/voicedock/extensions/tts/v1"
	"github.com/voicedock/ttspiper/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		fmt.Printf("failed to listen GRPC server: %s\n", err)
	}

	dataDir := "/data/dataset"
	dl := config.NewDownloader()
	cr := config.NewConfReader("/data/config/ttspiper.json")
	dr := config.NewDataReader(dataDir)
	cs := config.NewService(cr, dr, dl, dataDir)
	cs.LoadConfig()

	srv := grpcapi.NewServerTts(cs)

	s := grpc.NewServer()
	ttsv1.RegisterTtsAPIServer(s, srv)
	reflection.Register(s)
	s.Serve(lis)
}
