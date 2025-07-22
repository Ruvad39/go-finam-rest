package finam

import (
	"context"
	"net/http"
)

// Account Информация о конкретном аккаунте
type Account struct {
	AccountId        string  `json:"account_id,omitempty"`        // Идентификатор аккаунта
	Type             string  `json:"type,omitempty"`              // Тип аккаунта
	Status           string  `json:"status,omitempty"`            // Статус аккаунта
	Equity           Decimal `json:"equity,omitempty"`            // Доступные средства плюс стоимость открытых позиций
	UnrealizedProfit Decimal `json:"unrealized_profit,omitempty"` // Нереализованная прибыль
	Cash             []Money `json:"cash,omitempty"`              // Доступные средства
}

type AccountResponse struct {
	Account
	Positions []*Position `json:"positions,omitempty"` // Позиции. Открытые, плюс теоретические (по неисполненным активным заявкам)

}

// Информация о позиции
type Position struct {
	Symbol       string  `json:"symbol,omitempty"`        // Символ инструмента
	Quantity     Decimal `json:"quantity,omitempty"`      // Количество в шт., значение со знаком определяющее (long-short)
	AveragePrice Decimal `json:"average_price,omitempty"` // Средняя цена
	CurrentPrice Decimal `json:"current_price,omitempty"` // Текущая цена
}

// AccountRequest Получение Информация о конкретном аккаунте
type AccountRequest struct {
	client    *Client
	accountId string
}

func (c *Client) NewAccountRequest(accountId string) *AccountRequest {
	return &AccountRequest{
		client:    c,
		accountId: accountId,
	}
}

// Получение информации по конкретному аккаунту
// https://api.finam.ru/v1/accounts/account_id
//
// в запросе account_id - ваш номер счета
// в Headers - ваш jwt-token
func (r *AccountRequest) Do(ctx context.Context) (AccountResponse, error) {
	var err error
	var result AccountResponse
	req := NewRequest(http.MethodGet, apiURL).URLJoin("v1/accounts").URLJoin(r.accountId)
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
