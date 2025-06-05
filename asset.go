package finam

import (
	"context"
	"fmt"
	"net/http"
)

type Asset struct {
	Symbol string `json:"symbol,omitempty"` //  Символ инструмента ticker@mic
	Id     string `json:"id,omitempty"`
	Ticker string `json:"ticker,omitempty"` // Тикер инструмента
	Mic    string `json:"mic,omitempty"`    // mic идентификатор биржи
	Isin   string `json:"isin,omitempty"`   // Isin идентификатор инструмента
	Type   string `json:"type,omitempty"`   // Тип инструмента
	Name   string `json:"name,omitempty"`   // Наименование инструмента
}

type AssetsResponse struct {
	Assets []Asset
}

// AssetsRequest Получение списка доступных инструментов, их описание
type AssetsRequest struct {
	client *Client
}

func (c *Client) NewAssetsRequest() *AssetsRequest {
	return &AssetsRequest{
		client: c,
	}
}

// Получение списка доступных инструментов, их описание
// https://ftrr01.finam.ru/v1/assets
func (r *AssetsRequest) Do(ctx context.Context) (AssetsResponse, error) {
	var err error
	var result AssetsResponse
	req := NewRequest(http.MethodGet, apiURL).URLJoin("v1/assets")
	req.authorization = true
	// или можно самому добавим заголовок с авторизацией (accessToken)
	//r.client.WithAuthToken(req)

	resp, err := r.client.SendRequest(req)
	if err != nil {
		return result, err
	}
	err = resp.DecodeJSON(&result)
	if err != nil {
		return result, err
	}
	return result, nil

}

// Запрос торговых параметров инструмента
type AssetParamsRequest struct {
	client    *Client
	symbol    string // Символ инструмента
	accountId string // ID аккаунта для которого будут подбираться торговые параметры
}

func (c *Client) NewAssetParamsRequest(symbol, accountId string) *AssetParamsRequest {
	return &AssetParamsRequest{
		client:    c,
		symbol:    symbol,
		accountId: accountId,
	}
}

// Do Получение торговых параметров по инструменту
// GET /v1/assets/SBER@MISX/params?account_id=1440399
func (r *AssetParamsRequest) Do(ctx context.Context) error {
	var err error
	//var result AssetsResponse
	req := NewRequest(http.MethodGet, apiURL).URLJoin("v1/assets", r.symbol, "params")
	req.SetQuery("account_id", r.accountId)
	// добавим заголовок с авторизацией (accessToken)
	//r.client.WithAuthToken(req)
	// или можно проставить
	req.authorization = true
	//log.Debug("AssetParamsRequest.Do", "req.FullURL()", req.FullURL())
	resp, err := r.client.SendRequest(req)
	if err != nil {
		return err
	}
	fmt.Println("resp", resp)
	//log.Debug("AssetParamsRequest.Do", slog.Any("resp.Body", resp.Body))
	return err
}

// Запрос торговых параметров инструмента
type AssetInfoRequest struct {
	client *Client
	symbol string // Символ инструмента
}
