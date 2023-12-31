# TTS Piper
Piper based [VoiceDock TTS](https://github.com/voicedock/voicedock-specs/tree/main/proto/voicedock/core/tts/v1) implementation.

> Provides gRPC API for high quality text-to-speech (raw PCM stream) based on [Piper](https://github.com/rhasspy/piper) project.
> Provides download of new languages and voices via API.

# Usage
Run docker container:
```bash
docker run --rm \
  -v "$(pwd)/config:/data/config" \
  -v "$(pwd)/dataset:/data/dataset" \
  -p 9999:9999 \
  ghcr.io/voicedock/ttspiper:latest ttspiper
```

Show more options:
```bash
docker run --rm ghcr.io/voicedock/ttspiper ttspiper -h
```
```
Usage: ttspiper [--grpcaddr GRPCADDR] [--config CONFIG] [--datadir DATADIR] [--loglevel LOGLEVEL] [--logjson]

Options:
  --grpcaddr GRPCADDR    gRPC API host:port [default: 0.0.0.0:9999, env: GRPC_ADDR]
  --config CONFIG        configuration file for models [default: /data/config/ttspiper.json, env: CONFIG]
  --datadir DATADIR      dataset directory [default: /data/dataset, env: DATA_DIR]
  --loglevel LOGLEVEL    log level: debug, info, warn, error, dpanic, panic, fatal [default: info, env: LOG_LEVEL]
  --logjson              set to true to use JSON format [env: LOG_JSON]
  --help, -h             display this help and exit
```
## API
See implementation in [proto file](https://github.com/voicedock/voicedock-specs/blob/main/proto/voicedock/core/tts/v1/tts_api.proto).

## FAQ
### How to add a new language?
1. Find voice from [sample page](https://rhasspy.github.io/piper-samples/) or [Hugging Face](https://huggingface.co/rhasspy/piper-voices/tree/main).
2. Copy link to download `*.onnx` and `*.onnx.json` voice file
3. Add voice to [ttspiper.json](config%2Fttspiper.json) config:
   ```json
   {
     "lang": "lang_code",
     "speaker": "speaker_name",
     "download_onnx_url": "url *.onnx",
     "download_onnx_json_url": "url *.onnx.json",
     "license": "license text to accept"
   }
    ```

### How to use preloaded voices?
1. Add voice to [ttspiper.json](config%2Fttspiper.json) config (leave "download_onnx_url" and "download_onnx_json_url" blank to disable downloads).
2. [Download](https://rhasspy.github.io/piper-samples/) voice
3. Save `*.onnx` and `*.onnx.json` voice file to directory `dataset/{lang}/{speaker}/` (replace `{lang}` to language code and `{speaker}` to speaker name from configuration file  `ttspiper.json`)


## CONTRIBUTING
Lint proto files:
```bash
docker run --rm -w "/work" -v "$(pwd):/work" bufbuild/buf:latest lint internal/api/grpc/proto
```
Generate grpc interface:
```bash
docker run --rm -w "/work" -v "$(pwd):/work" ghcr.io/voicedock/protobuilder:1.0.0 generate internal/api/grpc/proto --template internal/api/grpc/proto/buf.gen.yaml
```
Manual build docker image:
```bash
docker build -t ghcr.io/voicedock/ttspiper:latest .
```

## Thanks
 * [Michael Hansen](https://github.com/synesthesiam) - TTS Piper uses the original [piper](https://github.com/rhasspy/piper) code and implements a shared library for c binding to go code