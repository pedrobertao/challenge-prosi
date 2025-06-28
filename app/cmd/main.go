package main

import (
	"context"
	"log"

	"github.com/pedrobertao/challenge-prosi/app/internal/config"
	"github.com/pedrobertao/challenge-prosi/app/internal/handlers"
	"github.com/pedrobertao/challenge-prosi/app/internal/routes"
	"github.com/pedrobertao/challenge-prosi/app/internal/storage"
)

func main() {
	cfg := config.Load()

	db, err := storage.Connect(cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close(context.Background())

	handler := handlers.New(db)
	app := routes.Setup(handler)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("error on server listener", err)
	}
}
