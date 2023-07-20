package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DataReader struct {
	dataDir string
}

func NewDataReader(dataDir string) *DataReader {
	return &DataReader{
		dataDir: dataDir,
	}
}

func (d *DataReader) ReadData(voice *VoiceConf) (*VoiceData, error) {
	ret := &VoiceData{
		VoiceConf: voice,
	}
	dataPath := filepath.Join(d.dataDir, voice.Lang, voice.Speaker)
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		return ret, nil
	}

	files, err := os.ReadDir(dataPath)
	if err != nil {
		return ret, fmt.Errorf("failed to read data directory: %w", err)
	}

	for _, v := range files {
		if strings.HasSuffix(v.Name(), ".onnx") {
			ret.OnnxPath = filepath.Join(dataPath, v.Name())
		}

		if strings.HasSuffix(v.Name(), ".onnx.json") {
			ret.OnnxJsonPath = filepath.Join(dataPath, v.Name())
		}
	}

	vSpec := &VoiceSpec{}
	f, err := os.Open(ret.OnnxJsonPath)
	if err != nil {
		return ret, fmt.Errorf("failed to open onnx.json: %w", err)
	}

	err = json.NewDecoder(f).Decode(&vSpec)
	if err != nil {
		return ret, fmt.Errorf("failed to read onnx.json: %w", err)
	}

	ret.VoiceSpec = vSpec

	return ret, nil
}
