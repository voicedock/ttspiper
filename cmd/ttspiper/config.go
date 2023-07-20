package main

type AppConfig struct {
	GrpcAddr         string `arg:"env:GRPC_ADDR" help:"gRPC API host:port" default:"0.0.0.0:9999"`
	Config           string `arg:"env:CONFIG" help:"configuration file for models" default:"/data/config/aichatllama.json"`
	DataDir          string `arg:"env:DATA_DIR" help:"dataset directory" default:"/data/dataset"`
	LogLevel         string `arg:"env:LOG_LEVEL" help:"log level: debug, info, warn, error, dpanic, panic, fatal" default:"info"`
	LogJson          bool   `arg:"env:LOG_JSON" help:"set to true to use JSON format"`
	LlamaGpuLayers   int    `arg:"env:LLAMA_GPU_LAYERS" help:"gpu layers for llama"`
	LlamaContextSize int    `arg:"env:LLAMA_CONTEXT_SIZE" help:"context size for llama" default:"1024"`
	LlamaThreads     int    `arg:"env:LLAMA_THREADS" help:"threads for llama (default MAX)"`
	LlamaTokens      int    `arg:"env:LLAMA_TOKENS" help:"sets number of tokens to generate for llama" default:"128"`
	LlamaDebug       bool   `arg:"env:LLAMA_DEBUG" help:"debug flag for llama"`
}
