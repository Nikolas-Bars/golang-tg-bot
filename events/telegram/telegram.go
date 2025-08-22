package telegram

import "golang-bot/clients/telegram"

type ProcessorStruct struct {
	tg *telegram.Client
	offset int

}

func New(client *telegram.Client, storage int) {

}