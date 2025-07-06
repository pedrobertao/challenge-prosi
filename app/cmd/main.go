package main

import (
	"context"
	"log"

	"github.com/pedrobertao/challenge-prosi/app/internal/config"
	"github.com/pedrobertao/challenge-prosi/app/internal/handlers"
	"github.com/pedrobertao/challenge-prosi/app/internal/routes"
	"github.com/pedrobertao/challenge-prosi/app/internal/storage"
	"github.com/pedrobertao/challenge-prosi/app/lib/logger"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	if err := logger.Setup(cfg.ENV); err != nil {
		log.Fatal("failed to init logger", err)
	}

	db, err := storage.Connect(cfg.MongoURI, cfg.DBName)
	if err != nil {
		logger.Fatal("failed to connect to database:", zap.Error(err))
	}
	defer db.Close(context.Background())

	handler := handlers.New(db)
	app := routes.Setup(handler)

	if err := app.Listen(":" + cfg.Port); err != nil {
		logger.Fatal("error on server listener", zap.Error(err))
	}
}
