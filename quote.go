package finam

import (
	"context"
	"net/http"
	"time"
)

// Информация о котировке
type Quote struct {
	Symbol    string    `json:"symbol,omitempty"`    // Символ инструмента
	Timestamp time.Time `json:"timestamp,omitempty"` // Метка времени
	Ask       Decimal   `json:"ask,omitempty"`       // Аск. 0 при отсутствии активного аска
	AskSize   Decimal   `json:"askSiz,omitempty"`    // Размер аска
	Bid       Decimal   `json:"bid,omitempty"`       // Бид. 0 при отсутствии активного бида
	BidSize   Decimal   `json:"bidSize,omitempty"`   // Размер бида
	Last      Decimal   `json:"last,omitempty"`      // Цена последней сделки
	LastSize  Decimal   `json:"lastSize,omitempty"`  // Размер последней сделки
	Volume    Decimal   `json:"volume,omitempty"`    // Дневной объем сделок
	Turnover  Decimal   `json:"turnover,omitempty"`  // Дневной оборот сделок
	Open      Decimal   `json:"open,omitempty"`      // Цена открытия. Дневная
	High      Decimal   `json:"high,omitempty"`      // Максимальная цена. Дневная
	Low       Decimal   `json:"low,omitempty"`       // Минимальная цена. Дневная
	Close     Decimal   `json:"close,omitempty"`     // Цена закрытия. Дневная
	Change    Decimal   `json:"change,omitempty"`    // Изменение цены (last минус close)
}

type QuoteRequest struct {
	client *Client
	symbol string
}

type QuoteResponse struct {
	Symbol string `json:"symbol,omitempty"` // Символ инструмента
	Quote  Quote  `json:"quote,omitempty"`  // Информация о котировке
}

func (c *Client) NewQuoteRequest(symbol string) *QuoteRequest {
	return &QuoteRequest{
		client: c,
		symbol: symbol,
	}
}

// Получение последней котировки по инструменту
// https://api.finam.ru/v1/instruments/YDEX@MISX/quotes/latest
func (r *QuoteRequest) Do(ctx context.Context) (QuoteResponse, error) {
	var err error
	var result QuoteResponse
	req := NewRequest(http.MethodGet, apiURL).URLJoin("v1/instruments").URLJoin(r.symbol).URLJoin("quotes/latest")
	req.authorization = true
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
