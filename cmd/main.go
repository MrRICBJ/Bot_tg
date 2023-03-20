package main

import (
	cli "bot/internal/clients/telegram"
	"bot/internal/config"
	"bot/internal/consumer/event-consumer"
	"bot/internal/events/telegram"
	"bot/internal/storage/database"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	tgBotHost = "api.telegram.org"
	batchSize = 100
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("error initializing configs: %s\n", err.Error())
	}

	cfg := config.MustLoad()

	ctx, cancel := context.WithCancel(context.Background())

	db := db.New(ctx, cfg)
	defer db.Db.Close()

	eventsProcessor := telegram.New(
		cli.New(tgBotHost, cfg.TgBotToken),
		db,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	go func() {
		if err := consumer.Run(ctx); err != nil {
			log.Fatal("service is stopped", err)
		}
	}()
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigint
		log.Print("stopping service...")
		time.Sleep(3 * time.Second)
		cancel()
	}()

	<-ctx.Done()
}
