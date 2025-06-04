package finam

import (
	"context"
	"net/http"
	"time"
)

// Информация о заявке
type Order struct {
	AccountId     string    `json:"accountId,omitempty"`     // Идентификатор аккаунта
	Symbol        string    `json:"symbol,omitempty"`        // Символ инструмента
	Quantity      *Decimal  `json:"quantity,omitempty"`      // Количество в шт.
	Side          Side      `json:"side,omitempty"`          // Сторона (long или short)
	Type          OrderType `json:"type,omitempty"`          // Тип заявки
	TimeInForce   string    `json:"timeInForce,omitempty"`   // Срок действия заявки
	LimitPrice    *Decimal  `json:"limitPrice,omitempty"`    // Необходимо для лимитной и стоп лимитной заявки
	StopPrice     *Decimal  `json:"stopPrice,omitempty"`     // Необходимо для стоп рыночной и стоп лимитной заявки
	StopCondition string    `json:"stopCondition,omitempty"` // Необходимо для стоп рыночной и стоп лимитной заявки
	ClientOrderId string    `json:"clientOrderId,omitempty"` // Уникальный идентификатор заявки. Автоматически генерируется, если не отправлен. (максимум 20 символов)
}

// Состояние заявки
type OrderState struct {
	// Идентификатор заявки
	OrderId string `json:"orderId,omitempty"`
	// Идентификатор исполнения
	ExecId string `json:"execId,omitempty"`
	// Статус заявки
	Status string `json:"status,omitempty"`
	// Заявка
	Order *Order `json:"order,omitempty"`
	// Дата и время выставления заявки
	TransactAt time.Time `json:"transactAt,omitempty"`
	//TransactAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=transact_at,json=transactAt,proto3" json:"transact_at,omitempty"`
	// Дата и время принятия заявки
	AcceptAt time.Time `json:"acceptAt,omitempty"`
	//AcceptAt *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=accept_at,json=acceptAt,proto3" json:"accept_at,omitempty"`
	// Дата и время  отмены заявки
	WithdrawAt time.Time `json:"withdrawAt,omitempty"`
	//WithdrawAt    *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=withdraw_at,json=withdrawAt,proto3" json:"withdraw_at,omitempty"`
}

// Список активных торговых заявок
type OrdersResponse struct {
	// Заявки
	Orders []*OrderState `protobuf:"bytes,1,rep,name=orders,proto3" json:"orders,omitempty"`
}

// GetOrdersRequest Получение списка ордеров по счету
type GetOrdersRequest struct {
	client    *Client
	accountId string
}

func (c *Client) NewGetOrdersRequest(accountId string) *GetOrdersRequest {
	return &GetOrdersRequest{
		client:    c,
		accountId: accountId,
	}
}

// Получение списка ордеров
// hhttps://api.finam.ru/v1/accounts/account_id/orders
//
// в запросе account_id - ваш номер счета
// в Headers - ваш jwt-token
func (r *GetOrdersRequest) Do(ctx context.Context) (OrdersResponse, error) {
	var err error
	var result OrdersResponse
	req := NewRequest(http.MethodGet, apiURL).URLJoin("v1/accounts").URLJoin(r.accountId).URLJoin("orders")
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
	//log.Info("OrdersRequest.Do", slog.Any("resp", resp))
	return result, nil

}

// CancelOrderRequest Запрос отмены торговой заявки
type CancelOrderRequest struct {
	client    *Client
	accountId string
	orderId   string
}

func (c *Client) NewCancelOrderRequest(accountId, orderId string) *CancelOrderRequest {
	return &CancelOrderRequest{
		client:    c,
		accountId: accountId,
		orderId:   orderId,
	}
}

// Отмена биржевой заявки DELETE
// https://api.finam.ru/v1/accounts/{account_id}/orders/{order_id}
//
// в запросе account_id - ваш номер счета
// в запросе order_id - идентификатор заявки
// в Headers - ваш jwt-token
func (r *CancelOrderRequest) Do(ctx context.Context) (OrderState, error) {
	var err error
	var result OrderState
	req := NewRequest(http.MethodDelete, apiURL)
	req.URLJoin("v1/accounts").URLJoin(r.accountId).URLJoin("orders").URLJoin(r.orderId)
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
	//log.Info("CancelOrderRequest.Do", slog.Any("resp", resp))
	return result, nil

}

// CancelAllOrders отменить все лимитные ордера
func (c *Client) CancelAllOrders(ctx context.Context, accountId string) error {
	// 1 список ордеров
	orders, err := c.NewGetOrdersRequest(accountId).Do(ctx)
	if err != nil {
		return err
	}
	for n, row := range orders.Orders {
		log.Debug("list Orders", "n", n,
			"state", row,
			"order", row.Order)
		id := row.OrderId
		state, errCancel := c.NewCancelOrderRequest(accountId, id).Do(ctx)
		if errCancel != nil {
			log.Error("CancelAllOrders", "err", errCancel.Error())
		}
		log.Debug("CancelAllOrders", "state", state)
	}

	return nil

}

// PlaceOrderRequest Запрос создания новой заявки
type PlaceOrderRequest struct {
	client *Client
	order  *Order
}

func (c *Client) NewPlaceOrderRequest() *PlaceOrderRequest {
	order := Order{
		Type:        OrderTypeMarket,
		TimeInForce: "TIME_IN_FORCE_DAY",
	}
	return &PlaceOrderRequest{
		client: c,
		order:  &order,
	}
}

// AccountId установить счет
func (r *PlaceOrderRequest) AccountId(value string) *PlaceOrderRequest {
	r.order.AccountId = value
	return r
}

// Symbol установить символ
func (r *PlaceOrderRequest) Symbol(value string) *PlaceOrderRequest {
	r.order.Symbol = value
	return r
}

func (r *PlaceOrderRequest) Quantity(value int) *PlaceOrderRequest {
	r.order.Quantity = IntToDecimal(value)
	return r
}

// Buy покупка по рынку
func (r *PlaceOrderRequest) Buy() *PlaceOrderRequest {
	r.order.Type = OrderTypeMarket
	r.order.Side = SideTypeBuy
	return r
}

// BuyLimit покупка по лимитной цене
func (r *PlaceOrderRequest) BuyLimit() *PlaceOrderRequest {
	r.order.Type = OrderTypeLimit
	r.order.Side = SideTypeBuy
	return r
}

// Sell продажа по рынку
func (r *PlaceOrderRequest) Sell() *PlaceOrderRequest {
	r.order.Type = OrderTypeMarket
	r.order.Side = SideTypeSell
	return r
}

// Sell продажа по лимитной цене
func (r *PlaceOrderRequest) SellLimit() *PlaceOrderRequest {
	r.order.Type = OrderTypeLimit
	r.order.Side = SideTypeSell
	return r
}

func (r *PlaceOrderRequest) Side(side Side) *PlaceOrderRequest {
	r.order.Side = side
	return r
}
func (r *PlaceOrderRequest) Type(orderType OrderType) *PlaceOrderRequest {
	r.order.Type = orderType
	return r
}

// LimitPrice установить цену для лимитной и стоп лимитной заявки
func (r *PlaceOrderRequest) LimitPrice(price float64) *PlaceOrderRequest {
	r.order.LimitPrice = Float64ToDecimal(price)
	return r
}

// StopPrice установить цену для стоп рыночной и стоп лимитной заявки
func (r *PlaceOrderRequest) StopPrice(price float64) *PlaceOrderRequest {
	r.order.StopPrice = Float64ToDecimal(price)
	return r
}

// Order
func (r *PlaceOrderRequest) Order(value *Order) *PlaceOrderRequest {
	r.order = value
	return r
}

// PlaceOrder
// POST https://api.finam.ru/v1/accounts/account_id/orders
// в запросе account_id - ваш номер счета
// в Headers - ваш jwt-token
// в body raw json Order
func (r *PlaceOrderRequest) Do(ctx context.Context) (OrderState, error) {
	var err error
	var result OrderState
	req := NewRequest(http.MethodPost, apiURL)
	req.URLJoin("v1/accounts").URLJoin(r.order.AccountId).URLJoin("orders")
	req.SetJSONBody(r.order) // в body raw json Order
	req.authorization = true
	resp, err := r.client.SendRequest(req)
	if err != nil {
		return result, err
	}
	err = resp.DecodeJSON(&result)
	if err != nil {
		return result, err
	}
	//log.Info("PlaceOrderRequest.Do", slog.Any("resp", resp))
	return result, nil
}
