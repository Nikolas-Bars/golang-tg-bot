package telegram

// Импортируем необходимые пакеты для работы с HTTP-запросами,
// обработки JSON и работы с URL
import (
	"encoding/json"
	"golang-bot/lib/e"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

// Client - основная структура для работы с Telegram API
// Содержит настройки подключения и HTTP-клиент
type Client struct {
	host     string      // адрес API Telegram
	basePath string      // базовый путь для запросов, включает токен бота
	client   http.Client // HTTP-клиент для отправки запросов
}

// Константы с названиями методов Telegram API
const getUpdateMethod = "getUpdates"    // метод для получения обновлений
const sendMessageMethod = "sendMessage" // метод для отправки сообщений

// New создает новый экземпляр клиента Telegram API
// Принимает токен бота и адрес API
func New(token string, host string) Client {
	return Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

// newBasePath формирует базовый путь для API запросов
// добавляя префикс "bot" к токену
func newBasePath(token string) string {
	return "bot" + token
}

// SendMessage отправляет сообщение в указанный чат
func (c Client) SendMessage(chatID int, text string) error {
	// создаем параметры запроса
	q := url.Values{} 
	// Функция strconv.Itoa в Go используется для преобразования целого числа (integer) в строку
	q.Add("chat_id", strconv.Itoa(chatID))
	// добавляем текст сообщения
	q.Add("text", text)                   

	_, err := c.doRequest(sendMessageMethod, q)

	if err != nil {
		return e.Wrap("can`t send message", err)
	}

	return nil
}

// Updates получает список обновлений (новых сообщений) от Telegram API
func (c Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	// конвертируем int в текст для параметра запроса
	q.Add(
		"offset",
		strconv.Itoa(offset),
	)
	// добавляем ограничение на количество обновлений
	q.Add("limit", strconv.Itoa(limit)) 

	// Выполняем запрос к API для получения обновлений
	data, err := c.doRequest(getUpdateMethod, q)

	if err != nil {
		return nil, e.WrapIfErr("can`t do request:", err)
	}

	var res UpdateResponse
	// Преобразуем полученный JSON-ответ в структуру UpdateResponse
	// Первый аргумент - JSON данные, второй - указатель на структуру для заполнения
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, e.WrapIfErr("can`t do request:", err)
	}

	return res.Result, nil // возвращаем список обновлений
}

// doRequest выполняет HTTP-запрос к Telegram API
// method - название метода API
// query - параметры запроса
func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	// Отложенная функция для обработки ошибок
	defer func() { err = e.WrapIfErr("can`t do request:", err) }()

	// Формируем URL для запроса
	u := url.URL{
		Scheme: "https",                       // используем защищенное соединение
		Host:   c.host,                        // хост API
		Path:   path.Join(c.basePath, method), // путь формируется из базового пути и метода
	}

	// Создаем новый GET-запрос
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, err
	}

	// Добавляем параметры запроса в URL
	req.URL.RawQuery = query.Encode()

	// Выполняем запрос
	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	// Гарантируем закрытие тела ответа после завершения функции
	defer func() { _ = resp.Body.Close() }()

	// Читаем тело ответа полностью
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	// возвращаем полученные данные
	return body, nil 
}
