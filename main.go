package main

import (
	"flag"
	"fmt"
	tgClient "golang-bot/clients/telegram"
	eventconsumer "golang-bot/consumer/event-consumer"
	"golang-bot/events/telegram"
	"golang-bot/lib/e/storage/files"
	"log"
)

const tgBotHost = "api.telegram.org"
const batchSize int = 100

const storagePath string = "storage"

func main() {
	
	token, host := mustToken()

	fmt.Printf("Token: %v \n", token)

	storage := files.New(storagePath)

	eventProcessor := telegram.New(
		tgClient.New(token, host),
		&storage,
	)

	consumer := eventconsumer.New(eventProcessor, eventProcessor, batchSize);

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stoped", err)
	}
	fmt.Printf("%v", eventProcessor)
}

func mustToken() (string, string) {
	token := flag.String(
		"token-bot-token",
		"",
		"token for access",
	)

	host := flag.String(
		"host-bot",
		"",
		"token for access",
	)

	flag.Parse()

	if *token == "" {
		// под капотом будет os.Exit(1)
		log.Fatal("token no valid")
	}
	if *host == "" {
		return *token, tgBotHost
	}
	return *token, *host
}
