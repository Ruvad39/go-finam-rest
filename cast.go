package finam

import (
	"fmt"
	"google.golang.org/genproto/googleapis/type/decimal"
	"strconv"
)

// IntToDecimal конвертируем int в google.Decimal
func IntToDecimal(i int) *decimal.Decimal {
	return &decimal.Decimal{
		Value: strconv.FormatInt(int64(i), 10),
	}
}

// Float64ToDecimal конвертируем  float64 в google.Decimal
func Float64ToDecimal(f float64) *decimal.Decimal {
	// Конвертируем float64 в строку с нужной точностью (например, 6 знаков после точки)
	// Можно использовать fmt.Sprintf("%.Nf", f) для фиксации количества знаков
	return &decimal.Decimal{
		Value: strconv.FormatFloat(f, 'f', -1, 64),
	}
}

// DecimalToFloat64E конвертируем google.Decimal в float64
// с обработкой ошибки
func DecimalToFloat64E(d *decimal.Decimal) (float64, error) {
	if d == nil {
		return 0, fmt.Errorf("decimal is nil")
	}
	return strconv.ParseFloat(d.Value, 64)
}

// DecimalToFloat64 конвертируем google.Decimal в float64
// БЕЗ обработки ошибки
func DecimalToFloat64(d *decimal.Decimal) float64 {
	result, _ := DecimalToFloat64E(d)
	return result
}

func DecimalToIntE(d *decimal.Decimal) (int, error) {
	if d == nil {
		return 0, nil
	}
	val, err := DecimalToFloat64E(d)
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

func DecimalToInt(d *decimal.Decimal) int {
	result, _ := DecimalToIntE(d)
	return result
}

func valueOrZero(v *decimal.Decimal) string {
	if v == nil {
		return "0"
	}
	return v.Value
}
