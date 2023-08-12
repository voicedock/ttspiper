package config

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Downloader struct {
}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) Download(url, outPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed download: %w", err)
	}
	defer resp.Body.Close()

	return d.saveFile(resp.Body, outPath)
}

func (d *Downloader) saveFile(r io.Reader, outPath string) error {
	outDir := filepath.Dir(outPath)
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err = os.MkdirAll(outDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create out directory: %w", err)
		}
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	if _, err := io.Copy(outFile, r); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	outFile.Close()

	return nil
}
