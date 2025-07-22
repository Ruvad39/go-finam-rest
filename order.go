package finam

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Информация о заявке
type Order struct {
	AccountId     string    `json:"account_id,omitempty"`      // Идентификатор аккаунта
	Symbol        string    `json:"symbol,omitempty"`          // Символ инструмента
	Quantity      *Decimal  `json:"quantity,omitempty"`        // Количество в шт.
	Side          Side      `json:"side,omitempty"`            // Сторона (long или short)
	Type          OrderType `json:"type,omitempty"`            // Тип заявки
	TimeInForce   string    `json:"time_in_force,omitempty"`   // Срок действия заявки
	LimitPrice    *Decimal  `json:"limit_price,omitempty"`     // Необходимо для лимитной и стоп лимитной заявки
	StopPrice     *Decimal  `json:"stop_price,omitempty"`      // Необходимо для стоп рыночной и стоп лимитной заявки
	StopCondition string    `json:"stop_condition,omitempty"`  // Необходимо для стоп рыночной и стоп лимитной заявки
	ClientOrderId string    `json:"client_order_id,omitempty"` // Уникальный идентификатор заявки. Автоматически генерируется, если не отправлен. (максимум 20 символов)
}

func (o Order) String() string {
	var builder strings.Builder

	builder.WriteString("Order{\n")
	builder.WriteString(fmt.Sprintf("\tAccountId: %q,\n", o.AccountId))
	builder.WriteString(fmt.Sprintf("\tSymbol: %q,\n", o.Symbol))

	if o.Quantity != nil {
		builder.WriteString(fmt.Sprintf("\tQuantity: %s,\n", o.Quantity.Value))
	} else {
		builder.WriteString("\tQuantity: nil,\n")
	}

	builder.WriteString(fmt.Sprintf("\tSide: %s,\n", o.Side))
	builder.WriteString(fmt.Sprintf("\tType: %s,\n", o.Type))
	builder.WriteString(fmt.Sprintf("\tTimeInForce: %q,\n", o.TimeInForce))

	if o.LimitPrice != nil {
		builder.WriteString(fmt.Sprintf("\tLimitPrice: %s,\n", o.LimitPrice.Value))
	} else {
		builder.WriteString("\tLimitPrice: nil,\n")
	}

	if o.StopPrice != nil {
		builder.WriteString(fmt.Sprintf("\tStopPrice: %s,\n", o.StopPrice.Value))
	} else {
		builder.WriteString("\tStopPrice: nil,\n")
	}

	builder.WriteString(fmt.Sprintf("\tStopCondition: %q,\n", o.StopCondition))
	builder.WriteString(fmt.Sprintf("\tClientOrderId: %q,\n", o.ClientOrderId))
	builder.WriteString("}")

	return builder.String()
}

// Состояние заявки
type OrderState struct {
	// Идентификатор заявки
	OrderId string `json:"order_id,omitempty"`
	// Идентификатор исполнения
	ExecId string `json:"exec_id,omitempty"`
	// Статус заявки
	Status string `json:"status,omitempty"`
	// Дата и время выставления заявки
	TransactAt time.Time `json:"transact_at,omitempty"`
	// Дата и время принятия заявки
	AcceptAt time.Time `json:"accept_at,omitempty"`
	// Дата и время отмены заявки
	WithdrawAt time.Time `json:"withdraw_at,omitempty"`
	// Заявка
	Order *Order `json:"order,omitempty"`
}

func (os OrderState) String() string {
	var builder strings.Builder

	builder.WriteString("OrderState{\n")
	builder.WriteString(fmt.Sprintf("\tOrderId: %q,\n", os.OrderId))
	builder.WriteString(fmt.Sprintf("\tExecId: %q,\n", os.ExecId))
	builder.WriteString(fmt.Sprintf("\tStatus: %q,\n", os.Status))
	builder.WriteString(fmt.Sprintf("\tTransactAt: %s,\n", os.TransactAt.Format(time.RFC3339)))
	builder.WriteString(fmt.Sprintf("\tAcceptAt: %s,\n", os.AcceptAt.Format(time.RFC3339)))
	builder.WriteString(fmt.Sprintf("\tWithdrawAt: %s,\n", os.WithdrawAt.Format(time.RFC3339)))
	builder.WriteString("}")
	if os.Order != nil {
		orderStr := strings.Replace(os.Order.String(), "\n", "\n\t", -1)
		builder.WriteString(fmt.Sprintf("\tOrder: %s,\n", orderStr))
	} else {
		builder.WriteString("\tOrder: nil,\n")
	}
	return builder.String()
}

// Список активных торговых заявок
type OrdersResponse struct {
	// Заявки
	Orders []*OrderState `json:"orders,omitempty"`
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
	//fmt.Println("resp", resp)
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
