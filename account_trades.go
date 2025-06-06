package finam

import (
	"context"
	"net/http"
	"time"
)

// Информация о сделке
type AccountTrade struct {
	TradeId   string    `json:"tradeId,omitempty"`   // Идентификатор сделки
	Symbol    string    `json:"symbol,omitempty"`    // Символ инструмента
	Price     Decimal   `json:"price,omitempty"`     // Цена исполнения
	Size      Decimal   `json:"size,omitempty"`      // Размер в шт.
	Side      Side      `json:"side,omitempty"`      // Сторона сделки (long или short)
	Timestamp time.Time `json:"timestamp,omitempty"` // Метка времени
	OrderId   string    `json:"orderId,omitempty"`   // Идентификатор заявки
}

// AccountTradesRequest Получение истории по сделкам аккаунта
type AccountTradesRequest struct {
	client    *Client
	accountId string    // Идентификатор аккаунта
	limit     int32     // Лимит количества сделок
	startTime time.Time // Начало запрашиваемого периода
	endTime   time.Time // Окончание запрашиваемого периода
}

// История по сделкам
type AccountTradesResponse struct {
	Trades []AccountTrade `json:"trades,omitempty"`
}

func (c *Client) NewAccountTradesRequest(accountId string) *AccountTradesRequest {
	return &AccountTradesRequest{
		client:    c,
		accountId: accountId,
	}
}

func (r *AccountTradesRequest) Limit(value int32) *AccountTradesRequest {
	r.limit = value
	return r
}

func (r *AccountTradesRequest) StartTime(value time.Time) *AccountTradesRequest {
	r.startTime = value
	return r
}

func (r *AccountTradesRequest) EndTime(value time.Time) *AccountTradesRequest {
	r.endTime = value
	return r
}

// Получение истории по сделкам аккаунта
// https://api.finam.ru/v1/accounts/account_id/trades?interval.start_time=2025-01-01T00:00:00Z&interval.end_time=2025-03-15T00:00:00Z
//
// в запросе account_id - ваш номер счета
// в Headers - ваш jwt-token
func (r *AccountTradesRequest) Do(ctx context.Context) (AccountTradesResponse, error) {
	var err error
	var result AccountTradesResponse
	req := NewRequest(http.MethodGet, apiURL).URLJoin("v1/accounts").URLJoin(r.accountId).URLJoin("trades")
	req.SetQuery("limit", r.limit)
	// time.RFC3339 =
	req.SetQuery("interval.start_time", r.startTime.Format(time.RFC3339))
	req.SetQuery("interval.end_time", r.endTime.Format(time.RFC3339))

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
