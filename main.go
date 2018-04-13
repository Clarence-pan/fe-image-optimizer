package main

import (
	"flag"
	"log"
	"regexp"
)

var cfg *tOptimizerConfig

var optimizedFileRe = regexp.MustCompile(`\.optimized\.\w+`)

func main() {
	var err error
	var configFileName string
	// var inputFileName string

	flag.StringVar(&configFileName, "config", "config.json", "specify the configuration file(default is config.json)")

	flag.Parse()

	cfg, err = loadConfig(configFileName)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("config: %#v", cfg)
	log.Printf("args: %#v", flag.Args())

	for _, x := range flag.Args() {
		doOptimize(x)
	}

}
