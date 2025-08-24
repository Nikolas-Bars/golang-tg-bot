package telegram

import (
	"errors"
	"golang-bot/clients/telegram"
	"golang-bot/events"
	"golang-bot/lib/e"
	"golang-bot/lib/e/storage"
)

var ErrUnknownType = errors.New("unknown event type")

var ErrGetMeta = errors.New("can`t get meta")

type ProcessorStruct struct {
	tg *telegram.Client
	offset int
	storage storage.Storage
}

type Meta struct {
	ChatID int
	Username string
}

func New(client *telegram.Client, storage storage.Storage) *ProcessorStruct {
	return &ProcessorStruct{
		tg: client,
		offset: 0,
		storage: storage,
	}
}
// возвращаем slice event`ов
func (p *ProcessorStruct) Fetch(limit int) ([]events.Event, error) {
	update, err := p.tg.Updates(p.offset, limit)

	if err != nil {
		return nil, e.Wrap("can`t get evemts", err)
	}

	if len(update) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(update))

	for _, u := range update {
		res = append(res, event(u))
	}

	p.offset = update[len(update)].ID + 1

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID: upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchType(u telegram.Update) events.Type {
	if u.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(u telegram.Update) string {
	if u.Message == nil {
		return ""
	}
	return u.Message.Text
}

func (p *ProcessorStruct) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		p.processMessage(event)
	default:
		return e.Wrap("can`t process message", ErrUnknownType) 
	}
}

func (p *ProcessorStruct) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.WrapIfErr("can`t process message", err)
	}
}

func meta(event events.Event)(Meta, error) {
	res, ok := event.Meta.(Meta)

	if !ok {
		return Meta{}, e.Wrap("can`t get meta", ErrGetMeta)
	}

	return res, nil
}