package main

import (
	"flag"
	"fmt"
	"log"
	"golang-bot/clients/telegram"
)

const tgBotHost = "api.telegram.org"

func main() {
	
	token, host := mustToken()

	fmt.Printf("Token: %v \n", token)

	tgClient :=  telegram.New(token, host)

	fmt.Printf("%v", tgClient)
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
		// под капотом будет os.Exit(1)
		return *token, tgBotHost
	}
	return *token, *host
}
