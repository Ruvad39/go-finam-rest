/*
замена  google.golang.org/genproto/googleapis/type/decimal"
*/
package finam

import (
	"fmt"
	"strconv"
)

// Decimal
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
