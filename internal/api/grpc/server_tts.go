package grpc

import (
	"context"
	"errors"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	ttsv1 "github.com/voicedock/ttspiper/internal/api/grpc/gen/voicedock/extensions/tts/v1"
	"github.com/voicedock/ttspiper/internal/config"
	"os"
	"sync"
	"unsafe"
)

/*
#cgo CFLAGS: -I/usr/local/include/ttspiperlib
#cgo LDFLAGS: -L/usr/local/lib -lttspiperlib -Wl,-rpath=/usr/local/lib
#include "ttspiperlib.h"
*/
import "C"

var mu sync.Mutex
var index int
var fns = make(map[int]func(data *C.int16_t, len C.int))

//export textToAudioCb
func textToAudioCb(cbId C.int, audioBuf *C.int16_t, audioBufLen C.int) {
	fn := lookup(int(cbId))
	fn(audioBuf, audioBufLen)
}

func register(cb func(data *C.int16_t, len C.int)) int {
	mu.Lock()
	defer mu.Unlock()
	index++
	for fns[index] != nil {
		index++
	}
	fns[index] = cb
	return index
}

func lookup(i int) func(data *C.int16_t, len C.int) {
	mu.Lock()
	defer mu.Unlock()
	return fns[i]
}

func unregister(i int) {
	mu.Lock()
	defer mu.Unlock()
	delete(fns, i)
}

func NewServerTts(configService *config.Service) *ServerTts {
	return &ServerTts{
		configService: configService,
	}
}

type ServerTts struct {
	configService *config.Service
	ttsv1.UnimplementedTtsAPIServer
}

func (s *ServerTts) TextToSpeech(in *ttsv1.TextToSpeechRequest, srv ttsv1.TtsAPI_TextToSpeechServer) error {
	C.initialize()
	defer C.terminate()

	voiceConfig := s.configService.FindDownloaded(in.Lang, in.Speaker)
	if voiceConfig == nil {
		return errors.New("voice not found")
	}

	// TODO: delete after stable release
	fw, _ := os.Create("/dataset/demo.wav")
	audioFormat := 1
	bitDepth := 16
	sampleRate := voiceConfig.VoiceSpec.Audio.SampleRate
	enc := wav.NewEncoder(fw, sampleRate, bitDepth, 1, audioFormat)

	defer enc.Close()

	voice := C.loadVoice(C.CString(voiceConfig.OnnxPath), C.CString(voiceConfig.OnnxJsonPath), nil)

	i := register(func(data *C.int16_t, length C.int) {
		slice := (*[1 << 28]C.int16_t)(unsafe.Pointer(data))[:length:length]
		out := make([]int32, 0, length)
		for _, v := range slice {
			out = append(out, int32(v))
		}

		enc.Write(&audio.IntBuffer{
			Format: &audio.Format{
				NumChannels: 1,
				SampleRate:  sampleRate,
			},
			Data:           ConvertInts[int](out),
			SourceBitDepth: bitDepth,
		})

		srv.Send(&ttsv1.TextToSpeechResponse{
			RawPcm:     out,
			SampleRate: int32(sampleRate),
			BitDepth:   int32(bitDepth),
		})

	})
	C.textToAudio(voice, C.CString(in.Text), C.int(i))
	unregister(i)

	return nil
}

func (s *ServerTts) GetVoices(ctx context.Context, in *ttsv1.GetVoicesRequest) (*ttsv1.GetVoicesResponse, error) {
	var voices []*ttsv1.Voice
	for _, v := range s.configService.FindAll() {
		voices = append(voices, &ttsv1.Voice{
			Lang:         v.VoiceConf.Lang,
			Speaker:      v.VoiceConf.Speaker,
			Downloaded:   v.Downloaded(),
			Downloadable: v.Downloadable(),
			License:      v.VoiceConf.License,
		})
	}

	return &ttsv1.GetVoicesResponse{
		Voices: voices,
	}, nil
}

func (s *ServerTts) DownloadVoice(ctx context.Context, in *ttsv1.DownloadVoiceRequest) (*ttsv1.DownloadVoiceResponse, error) {
	err := s.configService.Download(in.Lang, in.Speaker)

	return &ttsv1.DownloadVoiceResponse{}, err
}

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

func ConvertInts[U, T Int](s []T) (out []U) {
	out = make([]U, len(s))
	for i := range s {
		out[i] = U(s[i])
	}
	return out
}
