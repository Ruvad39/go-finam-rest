package finam

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Request struct {
	method        string
	URL           *url.URL
	header        http.Header
	query         url.Values
	body          io.Reader
	authorization bool // Делать авторизацию или нет
}

func NewRequest(method string, baseUrl string) *Request {
	base, _ := url.Parse(baseUrl)
	r := &Request{
		method:        method,
		URL:           base,
		authorization: true,
	}
	return r
}

// URLJoin добавим элементы пути к базовому
func (r *Request) URLJoin(elem ...string) *Request {
	// Если переданы элементы, то соединяем их с существующим путём
	if len(elem) > 0 {
		r.URL.Path = path.Join(r.URL.Path, path.Join(elem...))
	}
	return r
}

// SetHeader добавим заголовок
func (r *Request) SetHeader(key string, value string) *Request {
	if r.header == nil {
		r.header = make(http.Header)
	}
	//r.Header.Add(key, value)
	r.header.Set(key, value)

	return r
}

// SetQuery установим key/value для строки запроса (to query string)
func (r *Request) SetQuery(key string, value interface{}) *Request {
	if r.query == nil {
		r.query = url.Values{}
	}
	r.query.Set(key, fmt.Sprintf("%v", value))
	return r
}

// SetBody устанавливает тело запроса
func (r *Request) SetBody(body []byte) *Request {
	r.body = bytes.NewReader(body)
	return r
}

// SetJSONBody сериализует структуру в JSON и устанавливает как тело запроса.
func (r *Request) SetJSONBody(v any) *Request {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Errorf("failed to marshal JSON: %w", err))
	}
	r.body = bytes.NewReader(data)

	if r.header == nil {
		r.header = make(http.Header)
	}
	r.header.Set("Content-Type", "application/json")
	return r
}

// NewHttpRequest создадим и вернем *http.Request
// вызываем http.NewRequestWithContext
func (r *Request) NewHttpRequest(ctx context.Context) (*http.Request, error) {
	fullUrl := r.FullURL()
	req, err := http.NewRequestWithContext(ctx, r.method, fullUrl, r.body)
	if err != nil {
		return nil, err
	}
	// Копируем заголовки
	for key, values := range r.header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return req, nil
}

// Raw вернем *http.Request
func (r *Request) Raw() *http.Request {
	req, _ := r.NewHttpRequest(context.Background())
	return req

}

// FullURL возвращает полный URL (с путём и query-параметрами)
func (r *Request) FullURL() string {
	if r.URL == nil {
		return ""
	}

	// Копируем URL, чтобы не мутировать оригинальный объект
	u := *r.URL

	// Объединяем query из поля r.query (если он не пустой)
	if len(r.query) > 0 {
		q := u.Query()
		for key, values := range r.query {
			for _, value := range values {
				q.Add(key, value)
			}
		}
		u.RawQuery = q.Encode()
	}

	return u.String()
}

// Reset очистим данные структуры
func (r *Request) Reset() {
	r.method = ""
	r.URL = nil
	r.header = make(http.Header)
	r.query = make(url.Values)
	r.body = nil
}
