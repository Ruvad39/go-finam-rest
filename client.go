package finam

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net"
	"net/http"
	"time"
)

const (
	LibraryName = "FINAM-API-REST GO"
	Version     = "0.0.1"
	DateVersion = "2025-05-05 10:15:37"
)

// EndPoints
const apiURL = "https://ftrr01.finam.ru"

// http client
const defaultHTTPTimeout = time.Second * 10

var dialer = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
}

var defaultTransport = &http.Transport{
	Proxy:               http.ProxyFromEnvironment,
	DialContext:         dialer.DialContext,
	MaxIdleConns:        100,
	MaxConnsPerHost:     100,
	MaxIdleConnsPerHost: 100,
	// TLSNextProto:          make(map[string]func(string, *tls.Conn) http.RoundTripper),
	ExpectContinueTimeout: 0,
	ForceAttemptHTTP2:     true,
	TLSClientConfig:       &tls.Config{},
}

var DefaultHttpClient = &http.Client{
	Timeout:   defaultHTTPTimeout,
	Transport: defaultTransport,
}

// Client define API client
type Client struct {
	refreshToken string    // Refresh токен пользователя
	accessToken  string    // JWT токен для дальнейшей авторизации
	ttlJWT       time.Time // Время завершения действия JWT токена
	HttpClient   *http.Client
}

func NewClient(ctx context.Context, token string) (*Client, error) {
	client := &Client{
		refreshToken: token,
		HttpClient:   DefaultHttpClient,
	}
	err := client.initialize()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// initialize начальная иницилизация
func (c *Client) initialize() error {
	// получим и проставим accessToken
	err := c.UpdateJWT()
	if err != nil {
		return err
	}
	return nil
}

// Do вызов http запроса.
func (c *Client) SendRequest(req *Request) (*Response, error) {
	var err error
	// если запрос должен проходить с авторизацией
	// проставим в заголовок запроса c.accessToken
	if req.authorization {
		err = c.WithAuthToken(req)
		if err != nil {
			log.Error("SendRequest", "SetJWT", err.Error())
			return nil, err
		}
	}
	//log.Debug("SendRequest", slog.Any("req.Request", req.Raw()))
	resp, err := c.HttpClient.Do(req.Raw())
	if err != nil {
		log.Error("c.HttpClient.Do", "err", err.Error())
		return nil, err
	}
	//log.Debug("SendRequest", slog.Any("resp", resp))
	// newResponse reads the response body and return a new Response object
	response, err := NewResponse(resp)
	if err != nil {
		return response, err
	}

	// обработаем ошибку api
	if response.IsError() {
		log.Debug("SendRequest", slog.Any("response", response))
		//log.Debug("SendRequest", slog.Any("response.StatusCode", response.StatusCode), "response.Body", response.String())
		apiErr := makeAPIError(response)
		// return response, &responseErr{Code: response.StatusCode, Body: response.Body}
		return response, apiErr
	}

	return response, nil
}

// GetJson выполним запрос (SendRequest) и распарсим тело ответа (DecodeJSO)
func (c *Client) GetJson(req *Request, result interface{}) error {
	resp, err := c.SendRequest(req)
	if err != nil {
		return err
	}
	// в теле ответа ждем json
	if err = resp.DecodeJSON(&result); err != nil {
		return err
	}
	return nil
}
