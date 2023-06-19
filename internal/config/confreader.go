package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Reader struct {
	path string
}

func NewReader(path string) *Reader {
	return &Reader{
		path: path,
	}
}

func (r *Reader) ReadConfig() ([]*Voice, error) {
	var ret []*Voice
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