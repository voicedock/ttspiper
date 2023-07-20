package main

import (
	"github.com/alexflint/go-arg"
	grpcapi "github.com/voicedock/ttspiper/internal/api/grpc"
	ttsv1 "github.com/voicedock/ttspiper/internal/api/grpc/gen/voicedock/core/tts/v1"
	"github.com/voicedock/ttspiper/internal/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

var cfg AppConfig
var logger *zap.Logger

func init() {
	arg.MustParse(&cfg)
	logger = initLogger(cfg.LogLevel, cfg.LogJson)
}

func main() {
	defer logger.Sync()

	logger.Info(
		"Starting TTS Piper",
		zap.String("data_dir", cfg.DataDir),
		zap.String("config", cfg.Config),
	)

	lis, err := net.Listen("tcp", cfg.GrpcAddr)
	if err != nil {
		logger.Fatal("failed to listen GRPC server", zap.Error(err))
	}

	dl := config.NewDownloader()
	cr := config.NewConfReader(cfg.Config)
	dr := config.NewDataReader(cfg.DataDir)
	cs := config.NewService(cr, dr, dl, cfg.DataDir)
	err = cs.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load configuration", zap.Error(err))
	}

	srv := grpcapi.NewServerTts(cs, logger)

	s := grpc.NewServer()
	ttsv1.RegisterTtsAPIServer(s, srv)
	reflection.Register(s)

	logger.Info("gRPC server listen", zap.String("addr", cfg.GrpcAddr))
	err = s.Serve(lis)
	if err != nil {
		logger.Fatal("gRPC server error", zap.Error(err))
	}
}
