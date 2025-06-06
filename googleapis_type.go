/*
замена  google.golang.org/genproto/googleapis/type/decimal"
замена  google.golang.org/genproto/googleapis/type/money
замена google.golang.org/genproto/googleapis/type/date

*/

package finam

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Decimal
// замена  google.golang.org/genproto/googleapis/type/decimal"
type Decimal struct {
	Value string `json:"value"`
}

func NewDecimal(value string) *Decimal {
	return &Decimal{Value: value}
}

// IntToDecimal конвертируем int в Decimal
func IntToDecimal(i int) *Decimal {
	return &Decimal{
		Value: strconv.FormatInt(int64(i), 10),
	}
}

func (d *Decimal) Float64E() (float64, error) {
	if d == nil {
		return 0, fmt.Errorf("decimal is nil")
	}
	return strconv.ParseFloat(d.Value, 64)
}

// DecimalToFloat64 конвертируем Decimal в float64
// БЕЗ обработки ошибки
func (d *Decimal) Float64() float64 {
	result, _ := d.Float64E()
	return result
}

func (d *Decimal) IntE() (int, error) {
	if d == nil {
		return 0, nil
	}
	val, err := d.Float64E()
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

func (d *Decimal) Int() int {
	result, _ := d.IntE()
	return result
}

// Float64ToDecimal преобразует float64 в строку и сохраняет в Decimal.Value
func Float64ToDecimal(f float64) *Decimal {
	return &Decimal{
		Value: strconv.FormatFloat(f, 'f', -1, 64),
	}

}

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

// замена google.golang.org/genproto/googleapis/type/date
type Date struct {
	// Year of the date. Must be from 1 to 9999, or 0 to specify a date without
	// a year.
	Year int32 `json:"year,omitempty"`
	// Month of a year. Must be from 1 to 12, or 0 to specify a year without a
	// month and day.
	Month int32 `json:"month,omitempty"`
	// Day of a month. Must be from 1 to 31 and valid for the year and month, or 0
	// to specify a year by itself or a year and month where the day isn't
	// significant.
	Day int32 `json:"day,omitempty"`
	// contains filtered or unexported fields
}

func (d *Date) ToTimeE() time.Time {
	// Устанавливаем значения по умолчанию
	year := d.Year
	if year == 0 {
		year = 1
	}
	month := time.Month(d.Month)
	if month == 0 {
		month = time.January
	}
	day := d.Day
	if day == 0 {
		day = 1
	}

	return time.Date(
		int(year),
		month,
		int(day),
		0, 0, 0, 0,
		time.UTC,
	)
}

func (d *Date) ToTime() time.Time {
	// Устанавливаем значения по умолчанию
	year := d.Year
	if year == 0 {
		year = 1
	}
	month := time.Month(d.Month)
	if month == 0 {
		month = time.January
	}
	day := d.Day
	if day == 0 {
		day = 1
	}

	return time.Date(
		int(year),
		month,
		int(day),
		0, 0, 0, 0,
		time.UTC,
	)
}
