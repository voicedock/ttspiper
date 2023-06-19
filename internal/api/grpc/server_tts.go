package grpc

import (
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	ttsv1 "github.com/test/test/internal/api/grpc/gen/test/tts/v1"
	"os"
	"sync"
	"unsafe"
)

/*
#cgo CFLAGS: -I/usr/local/include/ttssimplelib
#cgo LDFLAGS: -L/usr/local/lib -lttssimplelib -Wl,-rpath=/usr/local/lib
#include "ttssimplelib.h"
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

func NewServerTts() *ServerTts {
	return &ServerTts{}
}

type ServerTts struct {
	ttsv1.UnimplementedTtsAPIServer
}

func (s *ServerTts) Convert(in *ttsv1.ConvertRequest, srv ttsv1.TtsAPI_ConvertServer) error {
	C.initialize()
	defer C.terminate()
	modelPath := "/dataset/ru-irinia-medium.onnx"
	modelConfigPath := "/dataset/ru-irinia-medium.onnx.json"
	//modelPath := "/dataset/en-us-lessac-low.onnx"
	//modelConfigPath := "/dataset/en-us-lessac-low.onnx.json"

	fw, _ := os.Create("/dataset/demo.wav")
	audioFormat := 1
	bitDepth := 16
	sampleRate := 22050
	enc := wav.NewEncoder(fw, sampleRate, bitDepth, 1, audioFormat)

	defer enc.Close()

	voice := C.loadVoice(C.CString(modelPath), C.CString(modelConfigPath), nil)

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

		srv.Send(&ttsv1.ConvertResponse{
			Chunk: out,
		})

	})
	C.textToAudio(voice, C.CString(in.Text), C.int(i))
	unregister(i)

	return nil
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
