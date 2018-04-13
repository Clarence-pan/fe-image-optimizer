package main

import "github.com/pkg/errors"

type tZipOptimizer struct {
	tBaseOptimizer
}

func newZipOptimizer(file string) *tZipOptimizer {
	return &tZipOptimizer{
		tBaseOptimizer: tBaseOptimizer{
			file: file,
		},
	}
}

func (opt *tZipOptimizer) optimize() error {
	return errors.New("Not implemented!")
}
