package main

import (
	"github.com/harmannkibue/actsml-jobs-orchestrator/config"
	"github.com/harmannkibue/actsml-jobs-orchestrator/internal/app"
	"log"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
