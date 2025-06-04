package finam

import (
	"context"
	"net/http"
	"time"

	//"github.com/shopspring/decimal"
	"google.golang.org/genproto/googleapis/type/decimal"
)

// Информация о заявке
type Order struct {
	// Идентификатор аккаунта
	AccountId string `json:"accountId,omitempty"`
	// Символ инструмента
	Symbol string `json:"symbol,omitempty"`
	// Количество в шт.
	Quantity *decimal.Decimal `json:"quantity,omitempty"`
	// Сторона (long или short)
	Side Side `json:"side,omitempty"`
	// Тип заявки
	Type OrderType `json:"type,omitempty"`
	// Срок действия заявки
	TimeInForce string `json:"timeInForce,omitempty"`
	// Необходимо для лимитной и стоп лимитной заявки
	LimitPrice *decimal.Decimal `json:"limitPrice,omitempty"`
	// Необходимо для стоп рыночной и стоп лимитной заявки
	StopPrice *decimal.Decimal `json:"stopPrice,omitempty"`
	// Необходимо для стоп рыночной и стоп лимитной заявки
	StopCondition string `json:"stopCondition,omitempty"`
	// Уникальный идентификатор заявки. Автоматически генерируется, если не отправлен. (максимум 20 символов)
	ClientOrderId string `json:"clientOrderId,omitempty"`
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

func (c *Client) NewPlaceOrderRequest(accountId, symbol string, side Side, quantity int) *PlaceOrderRequest {
	order := Order{
		AccountId:   accountId,
		Symbol:      symbol,
		Side:        side,
		Quantity:    IntToDecimal(quantity),
		Type:        OrderTypeMarket,
		TimeInForce: "TIME_IN_FORCE_DAY",
	}
	return &PlaceOrderRequest{
		client: c,
		order:  &order,
	}
}

func (r *PlaceOrderRequest) Side(side Side) *PlaceOrderRequest {
	r.order.Side = side
	return r
}
func (r *PlaceOrderRequest) Type(orderType OrderType) *PlaceOrderRequest {
	r.order.Type = orderType
	return r
}

func (r *PlaceOrderRequest) Price(price float64) *PlaceOrderRequest {
	r.order.LimitPrice = Float64ToDecimal(price)
	return r
}

func (r *PlaceOrderRequest) Quantity(value int) *PlaceOrderRequest {
	r.order.Quantity = IntToDecimal(value)
	return r
}

// POST PlaceOrder
// https://api.finam.ru/v1/accounts/account_id/orders
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
	//log.Info("PlaceOrderRequest.Do", slog.Any("resp", resp))
	return result, nil
}
