package main

type AppConfig struct {
	GrpcAddr string `arg:"env:GRPC_ADDR" help:"gRPC API host:port" default:"0.0.0.0:9999"`
	Config   string `arg:"env:CONFIG" help:"configuration file for models" default:"/data/config/ttspiper.json"`
	DataDir  string `arg:"env:DATA_DIR" help:"dataset directory" default:"/data/dataset"`
	LogLevel string `arg:"env:LOG_LEVEL" help:"log level: debug, info, warn, error, dpanic, panic, fatal" default:"info"`
	LogJson  bool   `arg:"env:LOG_JSON" help:"set to true to use JSON format"`
}
