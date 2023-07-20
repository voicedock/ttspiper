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
  -p 9997:9999 \
  ghcr.io/voicedock/ttspiper:latest /usr/src/app/main
```
## API
See implementation in [proto file](https://github.com/voicedock/voicedock-specs/blob/main/proto/voicedock/core/tts/v1/tts_api.proto).

## FAQ
### How to add a new language?
1. Find voice from [sample page](https://rhasspy.github.io/piper-samples/)
2. Copy link to download `tar.gz` voice file
3. Add voice to [ttspiper.json](config%2Fttspiper.json) config:
   ```json
   {
     "lang": "lang_code",
     "speaker": "speaker_name",
     "download_url": "download_url",
     "license": "license text to accept"
   }
    ```

### How to use preloaded voices?
1. Add voice to [ttspiper.json](config%2Fttspiper.json) config (leave "download_url" blank to disable downloads).
2. [Download](https://rhasspy.github.io/piper-samples/)
3. Extract voice to directory `dataset/{lang}/{speaker}/` (replace `{lang}` to language code and `{speaker}` to speaker name from configuration file  `ttspiper.json`)


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