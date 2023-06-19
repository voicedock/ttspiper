package config

import (
	"archive/tar"
	"compress/gzip"
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

	return d.extractFile(resp.Body, outPath)
}

func (d *Downloader) extractFile(r io.Reader, outPath string) error {
	decompressed, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("failed extract gzip: %w", err)
	}

	tarReader := tar.NewReader(decompressed)

	err = os.MkdirAll(outPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed create out directory: %w", err)
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("failed untar file: %w", err)
		}

		fullPath := filepath.Join(outPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(fullPath, 0755); err != nil {
				return fmt.Errorf("failed create dir: %w", err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(fullPath)
			if err != nil {
				return fmt.Errorf("failed create file: %w", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("failed write file: %w", err)
			}
			outFile.Close()

		default:
			return fmt.Errorf("uknown flag: %s in file %s", string(header.Typeflag), header.Name)
		}
	}

	return nil
}
