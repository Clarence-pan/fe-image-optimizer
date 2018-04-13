package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type tOptimizerConfig struct {
	Jpeg tJpegConfig `json:"jpeg"`
	Png  tPngConfig  `json:"png"`

	Jpegoptim string `json:"jpegoptim"`
	Pngquant  string `json:"pngquant"`
}

type tJpegConfig struct {
	MaxQuality int `json:"max_quality"`
}

type tPngConfig struct {
	MinQuality int `json:"min_quality"`
	MaxQuality int `json:"max_quality"`
}

func loadConfig(configFileName string) (cfg *tOptimizerConfig, err error) {
	buf, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open config file")
	}

	err = json.Unmarshal(buf, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "invalid configuration file")
	}

	err = cfg.sanitizeAndValidate()
	if err != nil {
		return nil, err
	}

	return
}

func (cfg *tOptimizerConfig) sanitizeAndValidate() error {

	if cfg.Jpeg.MaxQuality > 100 || cfg.Jpeg.MaxQuality < 0 {
		return errors.New("jpeg.max_quality must be in 0~100")
	}

	if cfg.Jpeg.MaxQuality == 0 {
		cfg.Jpeg.MaxQuality = 70
	}

	if cfg.Png.MaxQuality > 100 || cfg.Png.MaxQuality < 0 {
		return errors.New("png.max_quality must be in 0~100")
	}

	if cfg.Png.MaxQuality == 0 {
		cfg.Png.MaxQuality = 100
	}

	if cfg.Png.MinQuality > 100 || cfg.Png.MinQuality < 0 {
		return errors.New("png.min_quality must be in 0~100")
	}

	if cfg.Png.MinQuality == 0 {
		cfg.Png.MinQuality = minInt(50, cfg.Png.MaxQuality)
	}

	if cfg.Png.MinQuality > cfg.Png.MaxQuality {
		return errors.New("png.min_quality must be less than or equal with png.max_quality")
	}

	execFileDir := filepath.Dir(os.Args[0])

	if cfg.Jpegoptim == "" {
		cfg.Jpegoptim = filepath.Join(execFileDir, "lib", "jpegoptim")
	} else if !filepath.IsAbs(cfg.Jpegoptim) {
		cfg.Jpegoptim = filepath.Join(execFileDir, cfg.Jpegoptim)
	}

	if cfg.Pngquant == "" {
		cfg.Pngquant = filepath.Join(execFileDir, "lib", "pngquant")
	} else if !filepath.IsAbs(cfg.Pngquant) {
		cfg.Pngquant = filepath.Join(execFileDir, cfg.Pngquant)
	}

	return nil
}

func minInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}
