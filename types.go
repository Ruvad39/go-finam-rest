package finam

// Side Сторона сделки
type Side string

const (
	SideTypeBuy  Side = "SIDE_BUY"
	SideTypeSell Side = "SIDE_SELL"
)

// OrderType Тип заявки
type OrderType string

const (
	OrderTypeLimit     OrderType = "ORDER_TYPE_LIMIT"      // Лимитная заявка
	OrderTypeMarket    OrderType = "ORDER_TYPE_MARKET"     // Рыночная заявка
	OrderTypeStop      OrderType = "ORDER_TYPE_STOP"       // Стоп-заявка
	OrderTypeStopLimit OrderType = "ORDER_TYPE_STOP_LIMIT" // Стоп-лимитная заявка
)

// TimeInForceType define time in force type of order
type TimeInForce string

//const (
//	OrderTypeLimit     OrderType = "TIME_IN_FORCE_DAY"     // Лимитная заявка
//	OrderTypeMarket    OrderType = "ORDER_TYPE_MARKET"     // Рыночная заявка
//	OrderTypeStop      OrderType = "ORDER_TYPE_STOP"       // Стоп-заявка
//	OrderTypeStopLimit OrderType = "ORDER_TYPE_STOP_LIMIT" // Стоп-лимитная заявка
//)

// Статус заявки
type OrderStatus string

const (
	OrderStatusNew      OrderStatus = "ORDER_STATUS_NEW"      // На исполнении
	OrderStatusFilled   OrderStatus = "ORDER_STATUS_FILLED"   // Полностъю исполнилась (выполнилась)
	OrderStatusCanceled OrderStatus = "ORDER_STATUS_CANCELED" // Отменена
	OrderStatusRejected OrderStatus = "ORDER_STATUS_REJECTED" // отклонена
)
