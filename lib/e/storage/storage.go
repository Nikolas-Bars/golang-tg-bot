package storage

import (
	"crypto/sha1"
	"fmt"
	"golang-bot/lib/e"
	"io"
)

type Storage interface {
	Save(p *Page)
	PickRandom(userName string) (*Page, error)
	Remove(p *Page)
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

// Hash создает уникальный идентификатор страницы на основе URL и имени пользователя
// Возвращает хеш в виде шестнадцатеричной строки и ошибку, если что-то пошло не так
func (p Page) Hash() (string, error) {
	// Создаем новый SHA-1 хеш-объект
	// SHA-1 - это криптографическая хеш-функция, которая создает 160-битный (20-байтовый) хеш.
	h := sha1.New()

	// Записываем URL страницы в хеш
	// Используем WriteString для эффективной записи строки в хеш
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can`t calculate hash", err)
	}

	// Добавляем имя пользователя в хеш
	// Это гарантирует уникальность хеша даже если разные пользователи сохраняют одну страницу
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can`t calculate hash", err)
	}

	// Получаем итоговый хеш и конвертируем его в шестнадцатеричную строку
	// %x в fmt.Sprintf форматирует байты в hex-строку
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
