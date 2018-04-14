package main

import (
	"flag"
	"log"
)

var cfg *tOptimizerConfig

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

	if len(flag.Args()) == 0 {
		showMainWindow()
		return
	}

	// log.Printf("config: %#v", cfg)
	// log.Printf("args: %#v", flag.Args())

	for _, x := range flag.Args() {
		doOptimize(x)
	}

	log.Printf("all done.")
}
