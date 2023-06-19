package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type ConfReader struct {
	path string
}

func NewConfReader(path string) *ConfReader {
	return &ConfReader{
		path: path,
	}
}

func (r *ConfReader) ReadConfig() ([]*VoiceConf, error) {
	var ret []*VoiceConf
	f, err := os.Open(r.path)
	if err != nil {
		return nil, fmt.Errorf("failed open voice config: %w", err)
	}

	err = json.NewDecoder(f).Decode(&ret)
	if err != nil {
		return nil, fmt.Errorf("failed read voice config: %w", err)
	}

	return ret, nil
}
