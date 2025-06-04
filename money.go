/*
аналог google.golang.org/genproto/googleapis/type/money
https://github.com/googleapis/go-genproto/blob/513f23925822/googleapis/type/money/money.pb.go
*/
package finam

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Money struct {
	// The three-letter currency code defined in ISO 4217.
	CurrencyCode string `json:"currencyCode,omitempty"`
	// The whole units of the amount.
	// For example if `currencyCode` is `"USD"`, then 1 unit is one US dollar.
	Units int64 `json:"units,omitempty"`
	// Number of nano (10^-9) units of the amount.
	// The value must be between -999,999,999 and +999,999,999 inclusive.
	// If `units` is positive, `nanos` must be positive or zero.
	// If `units` is zero, `nanos` can be positive, zero, or negative.
	// If `units` is negative, `nanos` must be negative or zero.
	// For example $-1.75 is represented as `units`=-1 and `nanos`=-750,000,000.
	Nanos int32 `json:"nanos,omitempty"`
}

func (m *Money) Float64() float64 {
	if m == nil {
		return 0
	}
	return float64(m.Units) + float64(m.Nanos)/1e9
}

func (m *Money) String() string {
	amount := float64(m.Units) + float64(m.Nanos)/1e9
	return fmt.Sprintf("%s=%.2f", m.CurrencyCode, amount)
}

// UnmarshalJSON реализует кастомный парсинг для Money
// учитывающий, что поле units может быть как числом (int64), так и строкой (string):
// chat.deepseek.com
func (m *Money) UnmarshalJSON(data []byte) error {
	//if len(data) == 0 || string(data) == "null" {
	//	return nil // Обработка null или пустых данных
	//}
	type Alias Money // Создаем псевдоним, чтобы избежать рекурсии
	aux := &struct {
		Units interface{} `json:"units"` // Принимаем units как interface{}
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Обрабатываем units, который может быть string или int64
	switch v := aux.Units.(type) {
	case string:
		units, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid units value: %v", v)
		}
		m.Units = units
	case float64: // JSON числа всегда парсятся как float64
		m.Units = int64(v)
	case nil:
		m.Units = 0
	default:
		return fmt.Errorf("unexpected type for units: %T", v)
	}

	return nil
}
