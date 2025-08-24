package telegram

import (
	"golang-bot/lib/e"
	"golang-bot/lib/e/storage"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd = "/rnd"
	HelpCmd = "/help"
	StartCmd = "/start"
)

func (p *ProcessorStruct) doCmd(text string, chatId int, userName string) error {
	text = strings.TrimSpace(text)

	// обрабатываем логи нашего бота
	log.Printf("got new command '%s' from '%s'", text, userName)

	if isAddCmd(text) {

	}


	//команды - сохранить страницу https://..., получить рандомную страницу /rnd,
	// помощь - /help, и команда /start
	switch text{
	case RndCmd:
	case HelpCmd:
	case StartCmd:
	default:
	}
}

func (p *ProcessorStruct) savePage(chatID int, pageUrl string, username string) (err error) {
	defer func() {err = e.WrapIfErr("can`t save page", err)}()

	page := &storage.Page{
		URL: pageUrl,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)

	if err != nil {
		return err
	}

	if isExists {
		return p.tg.SendMessage(chatID, "this command exists")
		
	}

}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	// ссылки будут валидными только при наличии префикса https://
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}