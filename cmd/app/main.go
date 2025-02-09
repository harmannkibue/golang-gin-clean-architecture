package main

import (
	"github.com/harmannkibue/spectabill_psp_connector_clean_architecture/config"
	"github.com/harmannkibue/spectabill_psp_connector_clean_architecture/internal/app"
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
