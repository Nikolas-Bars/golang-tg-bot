package telegram

import (
	"errors"
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
		return p.savePage(chatId, text, userName)
	}


	//команды - сохранить страницу https://..., получить рандомную страницу /rnd,
	// помощь - /help, и команда /start
	switch text{
	case RndCmd:
		return p.sendRandom(chatId, userName)
	case HelpCmd:
		return p.sendHelp(chatId)
	case StartCmd:
		return p.sendHello(chatId)
	default:
		return p.tg.SendMessage(chatId, msgUnknownCommand)
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
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err;
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err;
	}
	return nil;
}

func (p *ProcessorStruct) sendRandom(chatID int, username string) (err error) {
	defer func() {err = e.WrapIfErr("can`t send random fucnc", err)}()

	page, err := p.storage.PickRandom(username)

	if err != nil && !errors.Is(err, storage.ErrorSavePages) {
		return err;
	}

	if errors.Is(err, storage.ErrorSavePages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err;
	}
	// удаляем страницу если нашли и отдали (как буд-то бы это лишнее)
	return p.storage.Remove(page)
}

func (p *ProcessorStruct) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *ProcessorStruct) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, messageHello)
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	// ссылки будут валидными только при наличии префикса https://
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}