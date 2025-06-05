package finam

import (
	"context"
	"net/http"
	"time"
)

// Время на сервере
type ClockResponse struct {
	// Метка времени
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// GetTime текущее время сервера
func (c *Client) GetTime(ctx context.Context) (time.Time, error) {
	var err error
	var result ClockResponse
	req := NewRequest(http.MethodGet, apiURL).URLJoin("v1/assets/clock")
	req.authorization = true
	resp, err := c.SendRequest(req)
	if err != nil {
		return time.Now(), err
	}
	err = resp.DecodeJSON(&result)
	if err != nil {
		return time.Now(), err
	}
	return result.Timestamp.In(TzMoscow), nil
}
