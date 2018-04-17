package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type tOptimizerConfig struct {
	Version int         `json:"version"`
	Jpeg    tJpegConfig `json:"jpeg"`
	Png     tPngConfig  `json:"png"`

	Jpegoptim string `json:"jpegoptim"`
	Pngquant  string `json:"pngquant"`

	OptimizedSuffix string   `json:"optimized_suffix"`
	BlackList       []string `json:"black_list"`
}

type tJpegConfig struct {
	MaxQuality int `json:"max_quality"`
}

type tPngConfig struct {
	MinQuality int `json:"min_quality"`
	MaxQuality int `json:"max_quality"`
}

func loadConfigOrDefault(configFileName string) (cfg *tOptimizerConfig, err error) {
	buf, err := ioutil.ReadFile(configFileName)
	if err != nil {
		cfg = &tOptimizerConfig{}
	} else {
		err = json.Unmarshal(buf, &cfg)
		if err != nil {
			return nil, errors.Wrap(err, "invalid configuration file")
		}
	}

	err = cfg.sanitizeAndValidate()
	if err != nil {
		return nil, err
	}

	return
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
		cfg.Jpeg.MaxQuality = 80
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

	if cfg.OptimizedSuffix == "" {
		cfg.OptimizedSuffix = ".o"
	}

	if cfg.BlackList == nil || len(cfg.BlackList) == 0 {
		cfg.BlackList = []string{
			"__MACOSX",
			"thumbs.db",
		}
	}

	return nil
}

func (cfg *tOptimizerConfig) isOptimizedFile(file string) bool {
	n := len(file)
	z := n - len(cfg.OptimizedSuffix)
	if z < 0 {
		return false
	}

	if file[z:] == cfg.OptimizedSuffix {
		return true
	}

	i := strings.LastIndex(file, ".")
	if i < 0 {
		return false
	}

	j := i - len(cfg.OptimizedSuffix)
	return j > 0 && file[j:i] == cfg.OptimizedSuffix
}

func (cfg *tOptimizerConfig) isInBlackList(file string) bool {
	a := strings.Split(strings.Replace(file, string(os.PathSeparator), "/", -1), "/")

	for _, y := range a {
		for _, x := range cfg.BlackList {
			if x == y {
				return true
			}
		}
	}

	return false
}

func (cfg *tOptimizerConfig) saveTo(file string) error {
	jsonBuf, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, jsonBuf, 0666)
}
