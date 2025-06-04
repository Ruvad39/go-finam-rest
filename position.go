package finam

// Информация о позиции
type Position struct {
	Symbol       string   `json:"symbol,omitempty"`       // Символ инструмента
	Quantity     *Decimal `json:"quantity,omitempty"`     // Количество в шт., значение со знаком определяющее (long-short)
	AveragePrice *Decimal `json:"averagePrice,omitempty"` // Средняя цена
	CurrentPrice *Decimal `json:"currentPrice,omitempty"` // Текущая цена
}
