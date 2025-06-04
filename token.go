/*
Срок действия access токена составляет 15 минут.
Ограничим 12 минут

TODO нужно ли запустить метод обновления токена в отдельном потоке с таймером?
*/
package finam

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const jwtTokenTtl = 12 // Время жизни токена JWT в минутах

// Запрос авторизации
type AuthRequest struct {
	// API токен (secret key)
	Secret string `protobuf:"bytes,1,opt,name=secret,proto3" json:"secret,omitempty"`
}

// Информация об авторизации
type AuthResponse struct {
	// Полученный JWT-токен
	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

// Запрос информации о токене
type TokenDetailsRequest struct {
	// JWT-токен
	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

// Информация о доступе к рыночным данным
type MDPermission struct {
	// Уровень котировок
	QuoteLevel string `json:"quoteLevel,omitempty"`
	// Задержка в минутах
	DelayMinutes int32 `json:"delayMinutes,omitempty"`
}

// Информация о токене
type TokenDetailsResponse struct {
	// Дата и время создания
	CreatedAt string `json:"createdAt,omitempty"`
	// Дата и время экспирации
	ExpiresAt string `json:"expiresAt,omitempty"`
	// Информация о доступе к рыночным данным
	//MdPermissions []*MDPermission `json:"mdPermissions,omitempty"`
	// Идентификаторы аккаунтов
	AccountIds []string `json:"accountIds,omitempty"`
}

// WithAuthToken добавим в запрос токен авторизации
//
// При необходимости, предварительно получим токен (GetJWT)
// и запишем его в параметры клиента (c.accessToken)
func (c *Client) WithAuthToken(req *Request) error {
	// проверим наличие токена
	err := c.UpdateJWT()
	if err != nil {
		return err
	}
	// запишем его в запрос
	if req != nil {
		req.SetHeader("Authorization", c.accessToken)
	}

	return nil
}

// UpdateJWT
// если токен пустой или вышло его время => получим JWT и запишем его в параметры клиента
func (c *Client) UpdateJWT() error {
	if c.refreshToken == "" {
		c.accessToken = ""
		return fmt.Errorf("UpdateJWT: token пустой")
	}
	// если токен пустой или вышло его время
	if c.accessToken == "" || c.ttlJWT.Before(time.Now()) {
		log.Debug("UpdateJWT. токен пустой или вышло его время = получим новый токен")
		// получим новый токен
		token, err := c.GetJWT()
		if err != nil {
			return err
		}
		// запишем время окончания токена
		c.ttlJWT = time.Now().Add(jwtTokenTtl * time.Minute)
		c.accessToken = token
	}
	log.Debug("UpdateJWT. токен живой")
	return nil
}

// GetJWT получим accessToken
//
// v1/sessions
func (c *Client) GetJWT() (string, error) {
	//const op = "GetJWT"
	var err error
	var result AuthResponse
	if c.refreshToken == "" {
		c.accessToken = ""
		return c.accessToken, err
	}
	req := NewRequest(http.MethodPost, apiURL).URLJoin("v1/sessions")
	reqAuth := &AuthRequest{
		Secret: c.refreshToken,
	}
	req.SetJSONBody(reqAuth) // запишем в тело запроса структуру
	req.authorization = false
	log.Debug("GetJWT start refresh")
	t := time.Now()
	err = c.GetJson(req, &result)
	if err != nil {
		return "", err
	}
	log.Debug("GetJWT end refresh", "duration", time.Since(t))
	return result.Token, nil

}

// GetTokenDetails Получение информации о токене сессии
//
// v1/sessions/details/
func (c *Client) GetTokenDetails(ctx context.Context) (TokenDetailsResponse, error) {
	var err error
	var result TokenDetailsResponse
	req := NewRequest(http.MethodPost, apiURL).URLJoin("v1/sessions/details/")
	reqToken := &TokenDetailsRequest{
		Token: c.accessToken,
	}
	req.SetJSONBody(reqToken)
	resp, err := c.SendRequest(req)
	if err != nil {
		return result, err
	}
	err = resp.DecodeJSON(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// WithAuthToken добавим в запрос токен авторизации
// func (c *Client) WithAuthToken(req *Request) error {
// 	return c.SetJWT(req)
// }
