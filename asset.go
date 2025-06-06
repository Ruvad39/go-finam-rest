package finam

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"github.com/valyala/fastjson"
	"math"
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

// информации по конкретному инструменту
type AssetInfo struct {
	Board          string `json:"board,omitempty"`          // Код режима торгов
	Id             string `json:"id,omitempty"`             // Идентификатор инструмента
	Ticker         string `json:"ticker,omitempty"`         // Тикер инструмента
	Mic            string `json:"mic,omitempty"`            // mic идентификатор биржи
	Isin           string `json:"isin,omitempty"`           // Isin идентификатор инструмента
	Type           string `json:"type,omitempty"`           // Тип инструмента
	Name           string `json:"name,omitempty"`           // Наименование инструмента
	Decimals       int32  `json:"decimals,omitempty"`       // Кол-во десятичных знаков в цене
	MinStep        int64  `json:"minStep,omitempty"`        // Минимальный шаг цены
	LotSize        int32  `json:"lotSize,omitempty"`        // Кол-во штук в лоте
	ExpirationDate Date   `json:"expirationDate,omitempty"` // Дата экспирации фьючерса
	//LotSize        Decimal `json:"lotSize,omitempty"`        // Кол-во штук в лоте
}

// StepPrice рассчитаем шаг цены
// шаг цены = float64(MinStep) * math.Pow(10, -float64(Decimals))
func (a *AssetInfo) StepPrice() float64 {
	return float64(a.MinStep) * math.Pow(10, -float64(a.Decimals))
}

// NormalizePrice приведем цену к точности шага цены инструмента
func (a *AssetInfo) NormalizePrice(price float64) float64 {
	stepPrice := float64(a.MinStep) * math.Pow(10, -float64(a.Decimals))
	if stepPrice != 0 {
		return float64(int64(price/stepPrice)) * stepPrice
	}
	return price
}

// UnmarshalJSON
// ручной парсинг JSON
// так как minStep приходит как символьное
// и сразу поменяю Decimal (lotSize) в int32
// TODO DATE сделать конвертацию в time?
func (a *AssetInfo) UnmarshalJSON(data []byte) error {
	var p fastjson.Parser
	v, err := p.ParseBytes(data)
	if err != nil {
		return err
	}
	// Обрабатываем каждое поле
	a.Board = string(v.GetStringBytes("board"))
	a.Id = string(v.GetStringBytes("id"))
	a.Ticker = string(v.GetStringBytes("ticker"))
	a.Mic = string(v.GetStringBytes("mic"))
	a.Isin = string(v.GetStringBytes("isin"))
	a.Type = string(v.GetStringBytes("type"))
	a.Name = string(v.GetStringBytes("name"))
	a.Decimals = int32(v.GetInt64("decimals"))
	// Особенная обработка minStep (может быть строкой или числом)
	minStepVal := string(v.GetStringBytes("minStep"))
	a.MinStep, err = cast.ToInt64E(minStepVal)
	if err != nil {
		return fmt.Errorf("invalid minStep value: %v", minStepVal)
	}

	// Обработка Decimal (lotSize)
	lotSizeVal := v.Get("lotSize")
	if lotSizeVal.Exists("value") {
		//a.LotSize.Value = string(lotSizeVal.GetStringBytes("value"))
		a.LotSize = cast.ToInt32(string(lotSizeVal.GetStringBytes("value")))
	}

	// Обработка Date (expirationDate)
	if v.Exists("expirationDate") {
		dateVal := v.Get("expirationDate")
		if dateVal.Exists("year") {
			a.ExpirationDate.Year = int32(dateVal.GetInt64("year"))
		}
		if dateVal.Exists("month") {
			a.ExpirationDate.Month = int32(dateVal.GetInt64("month"))
		}
		if dateVal.Exists("day") {
			a.ExpirationDate.Day = int32(dateVal.GetInt64("day"))
		}
	}
	return nil
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
	client    *Client
	symbol    string // Символ инструмента
	accountId string // ID аккаунта для которого будут подбираться торговые параметры
}

func (c *Client) NewAssetInfoRequest(symbol, accountId string) *AssetInfoRequest {
	return &AssetInfoRequest{
		client:    c,
		symbol:    symbol,
		accountId: accountId,
	}
}

// Do Получить информацию по конкретному инструменту
// GET /v1/assets/SBER@MISX?account_id=1440399
func (r *AssetInfoRequest) Do(ctx context.Context) (AssetInfo, error) {
	var err error
	var result AssetInfo
	req := NewRequest(http.MethodGet, apiURL).URLJoin("v1/assets", r.symbol)
	req.SetQuery("account_id", r.accountId)
	req.authorization = true
	resp, err := r.client.SendRequest(req)
	if err != nil {
		return result, err
	}
	// fmt.Println("resp", resp)
	err = resp.DecodeJSON(&result)
	if err != nil {
		return result, err
	}
	return result, err
}
