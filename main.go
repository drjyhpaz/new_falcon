package main

import (
	"falcon/config"
	"falcon/ui"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app := ui.NewApp(cfg)
	app.Run()
}
