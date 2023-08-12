package grpc

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	commonv1 "github.com/voicedock/ttspiper/internal/api/grpc/gen/voicedock/core/common/v1"
	ttsv1 "github.com/voicedock/ttspiper/internal/api/grpc/gen/voicedock/core/tts/v1"
	"github.com/voicedock/ttspiper/internal/config"
	"go.uber.org/zap"
	"sync"
	"unsafe"
)

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/src/app/lib -lttspiperlib -Wl,-rpath=/usr/src/app/lib
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

func NewServerTts(configService *config.Service, logger *zap.Logger) *ServerTts {
	return &ServerTts{
		configService: configService,
		logger:        logger,
	}
}

type ServerTts struct {
	configService *config.Service
	logger        *zap.Logger
	ttsv1.UnimplementedTtsAPIServer
}

func (s *ServerTts) TextToSpeech(in *ttsv1.TextToSpeechRequest, srv ttsv1.TtsAPI_TextToSpeechServer) error {
	s.logger.Info("TextToSpeech: starting")
	defer s.logger.Info("TextToSpeech: complete")
	C.initialize()
	defer C.terminate()

	voiceConfig := s.configService.FindDownloaded(in.Lang, in.Speaker)
	if voiceConfig == nil {
		return fmt.Errorf("voice not found by lang `%s` and speaker `%s`", in.Lang, in.Speaker)
	}

	sampleRate := voiceConfig.VoiceSpec.Audio.SampleRate

	voice := C.loadVoice(C.CString(voiceConfig.OnnxPath), C.CString(voiceConfig.OnnxJsonPath), nil)

	i := register(func(data *C.int16_t, length C.int) {
		slice := (*[1 << 28]C.int16_t)(unsafe.Pointer(data))[:length:length]
		out := make([]int16, 0, length)
		for _, v := range slice {
			out = append(out, int16(v))
		}

		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, out)

		err := srv.Send(&ttsv1.TextToSpeechResponse{
			Audio: &commonv1.AudioContainer{
				Data:       buf.Bytes(),
				SampleRate: int32(sampleRate),
				Channels:   1,
			},
		})
		if err != nil {
			s.logger.Error("TextToSpeech: failed to send audio", zap.Error(err))
		}
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
	s.logger.Info("DownloadVoice: starting",
		zap.String("lang", in.Lang), zap.String("speaker", in.Speaker))
	defer s.logger.Info("DownloadVoice: complete",
		zap.String("lang", in.Lang), zap.String("speaker", in.Speaker))

	err := s.configService.Download(in.Lang, in.Speaker)

	return &ttsv1.DownloadVoiceResponse{}, err
}
