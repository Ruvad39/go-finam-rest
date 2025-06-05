package finam

import (
	"net/http"
	"time"
)
import (
	"context"
)

// структуры данных для свечей

// Timeframe период свечей
type Timeframe string

const (
	TimeframeM1  Timeframe = "TIME_FRAME_M1"  // 1 минута. Глубина данных 7 дней.
	TimeframeM5  Timeframe = "TIME_FRAME_M5"  // 5 минут. Глубина данных 30 дней.
	TimeframeM15 Timeframe = "TIME_FRAME_M15" // 15 минут. Глубина данных 30 дней
	TimeframeM30 Timeframe = "TIME_FRAME_M30" // 30 минут. Глубина данных 30 дней.
	TimeframeH1  Timeframe = "TIME_FRAME_H1"  // 1 час. Глубина данных 30 дней.
	TimeframeH2  Timeframe = "TIME_FRAME_H2"  // 2 часа. Глубина данных 30 дней
	TimeframeH4  Timeframe = "TIME_FRAME_H4"  // 4 часа. Глубина данных 30 дней.
	TimeframeH8  Timeframe = "TIME_FRAME_H8"  // 8 часов. Глубина данных 30 дней.
	TimeframeD1  Timeframe = "TIME_FRAME_D"   // День. Глубина данных 365 дней
	TimeframeW1  Timeframe = "TIME_FRAME_W"   // Неделя. Глубина данных 365*5 дней.
)

// Информация об агрегированной свече
type Bar struct {
	Timestamp time.Time `json:"timestamp,omitempty"` // Метка времени
	Open      Decimal   `json:"open,omitempty"`      // Цена открытия свечи
	High      Decimal   `json:"high,omitempty"`      // Максимальная цена свечи
	Low       Decimal   `json:"low,omitempty"`       // Минимальная цена свечи
	Close     Decimal   `json:"close,omitempty"`     // Цена закрытия свечи
	Volume    Decimal   `json:"volume,omitempty"`    // Объём торгов за свечу
}

// Список агрегированных свеч
type BarsResponse struct {
	Timeframe Timeframe // период свечей
	Symbol    string    `json:"symbol,omitempty"` // Символ инструмента
	Bars      []Bar     `json:"bars,omitempty"`   // Агрегированная свеча
}

// Запрос торговых параметров инструмента
type BarsRequest struct {
	client    *Client
	symbol    string    // Символ инструмента
	timeframe Timeframe // Необходимый таймфрейм
	startTime time.Time // Начало запрашиваемого периода
	endTime   time.Time // Окончание запрашиваемого периода
}

func (c *Client) NewBarsRequest() *BarsRequest {
	return &BarsRequest{
		client: c,
	}
}
func (r *BarsRequest) Symbol(value string) *BarsRequest {
	r.symbol = value
	return r
}
func (r *BarsRequest) Timeframe(value Timeframe) *BarsRequest {
	r.timeframe = value
	return r
}
func (r *BarsRequest) StartTime(value time.Time) *BarsRequest {
	r.startTime = value
	return r
}

func (r *BarsRequest) EndTime(value time.Time) *BarsRequest {
	r.endTime = value
	return r
}

// Получение исторических данных по инструменту (агрегированные свечи)
// https://api.finam.ru/v1/instruments/YDEX@MISX/bars?interval.start_time=2025-03-01T00:00:00Z&interval.end_time=2025-03-15T00:00:00Z&timeframe=TIME_FRAME_D
func (r *BarsRequest) Do(ctx context.Context) (BarsResponse, error) {
	var err error
	var result BarsResponse
	result.Timeframe = r.timeframe

	req := NewRequest(http.MethodGet, apiURL).URLJoin("v1/instruments").URLJoin(r.symbol).URLJoin("bars")
	req.SetQuery("timeframe", r.timeframe)
	// time.RFC3339 =
	req.SetQuery("interval.start_time", r.startTime.Format(time.RFC3339))
	req.SetQuery("interval.end_time", r.endTime.Format(time.RFC3339))
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
