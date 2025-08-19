package telegram

import (
	"encoding/json"
	"golang-bot/lib/e"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const getUpdateMethod = "getUpdates"
const sendMessageMethod = "sendMessage"

func New(token string, host string) Client {
	return Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c Client) SendMessage(chatID int, text string) error {
	q := url.Values{}

	q.Add("chatID", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)

	if err != nil {
		return e.Wrap("can`t send message", err)
	}

	return nil
}

func (c Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add(
		"offset",
		strconv.Itoa(offset), // конвертируем int в текст
	)
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdateMethod, q)

	if err != nil {
		return nil, e.WrapIfErr("can`t do request:", err)
	}

	var res UpdateResponse
	// Unmarshal - распарсит json. 1 арг - данные, 2 арг - куда парсит (обязат указатель)
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, e.WrapIfErr("can`t do request:", err)
	}

	return res.Result, nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {

	defer func() { err = e.WrapIfErr("can`t do request:", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method), // метод сам уберет лишние слэши если они будут
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil

}
