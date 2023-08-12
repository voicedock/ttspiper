package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type Service struct {
	confReader *ConfReader
	dataReader *DataReader
	downloader *Downloader
	config     []*VoiceData
	idxConfig  map[string]map[string]*VoiceData
	dataDir    string
}

func NewService(
	confReader *ConfReader,
	dataReader *DataReader,
	downloader *Downloader,
	dataDir string,
) *Service {
	return &Service{
		confReader: confReader,
		dataReader: dataReader,
		downloader: downloader,
		config:     []*VoiceData{},
		idxConfig:  make(map[string]map[string]*VoiceData),
		dataDir:    dataDir,
	}
}

func (s *Service) LoadConfig() error {
	voiceConf, err := s.confReader.ReadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	var cfg []*VoiceData
	for _, v := range voiceConf {
		vData, _ := s.dataReader.ReadData(v)
		cfg = append(cfg, vData)
	}

	s.config = cfg

	s.RebuildIdx()
	return nil
}

func (s *Service) RebuildIdx() {
	for _, v := range s.config {
		if s.idxConfig[v.VoiceConf.Lang] == nil {
			s.idxConfig[v.VoiceConf.Lang] = make(map[string]*VoiceData)
		}

		s.idxConfig[v.VoiceConf.Lang][v.VoiceConf.Speaker] = v
	}
}

func (s *Service) FindAll() []*VoiceData {
	return s.config
}

func (s *Service) Download(lang, speaker string) error {
	voice, ok := s.idxConfig[lang][speaker]
	if !ok {
		return errors.New("voice not found")
	}

	if !voice.Downloadable() {
		return errors.New("voice is not downloadable")
	}

	outPath := filepath.Join(s.dataDir, lang, speaker)
	onnxName := filepath.Base(voice.VoiceConf.DownloadOnnxUrl)
	if !strings.HasSuffix(onnxName, ".onnx") {
		onnxName = "model.onnx"
	}
	jsonName := filepath.Base(voice.VoiceConf.DownloadOnnxJsonUrl)
	if !strings.HasSuffix(jsonName, ".onnx.json") {
		onnxName = "model.onnx.json"
	}

	err := s.downloader.Download(voice.VoiceConf.DownloadOnnxUrl, filepath.Join(outPath, onnxName))
	if err != nil {
		return fmt.Errorf("failed to download voice: %w", err)
	}
	err = s.downloader.Download(voice.VoiceConf.DownloadOnnxJsonUrl, filepath.Join(outPath, jsonName))
	if err != nil {
		return fmt.Errorf("failed to download voice: %w", err)
	}

	return s.LoadConfig()
}

func (s *Service) FindDownloaded(lang, speaker string) *VoiceData {
	ret := s.idxConfig[lang][speaker]
	if ret != nil && ret.Downloaded() {
		return ret
	}

	return nil
}
