package finam

import (
	"context"
	"net/http"
)

// Account Информация о конкретном аккаунте
type Account struct {
	AccountId        string  `json:"accountId,omitempty"`        // Идентификатор аккаунта
	Type             string  `json:"type,omitempty"`             // Тип аккаунта
	Status           string  `json:"status,omitempty"`           // Статус аккаунта
	Equity           Decimal `json:"equity,omitempty"`           // Доступные средства плюс стоимость открытых позиций
	UnrealizedProfit Decimal `json:"unrealizedProfit,omitempty"` // Нереализованная прибыль
	Cash             []Money `json:"cash,omitempty"`             // Доступные средства
}

type AccountResponse struct {
	Account
	Positions []*Position `json:"positions,omitempty"` // Позиции. Открытые, плюс теоретические (по неисполненным активным заявкам)

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
