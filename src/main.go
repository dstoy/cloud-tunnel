package main

import (
	"log"
	"os"

	"github.com/dstoy/tunnel/src/config"
	flag "github.com/spf13/pflag"
)

func main() {
	// Collect the application arguments
	var configFile *string = flag.String("config", "", "configuration file")
	flag.Parse()

	// Initialize the configuration
	err := config.Load(*configFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
