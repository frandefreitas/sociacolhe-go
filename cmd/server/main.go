package main

import (
	"log"
	"os"

	"sociacolhe/internal/config"
	"sociacolhe/internal/db"
	"sociacolhe/internal/routes"
)

func main() {
	cfg := config.Load()
	if err := db.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("db connect: %v", err)
	}

	r := routes.SetupRouter(cfg)
	log.Printf("listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
