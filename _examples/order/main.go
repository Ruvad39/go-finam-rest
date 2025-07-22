package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Ruvad39/go-finam-rest"
	"github.com/joho/godotenv"
)

func main() {
	// предполагаем что есть файл .env в котором записан secret-Token в переменной FINAM_TOKEN
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
	}
	token, _ := os.LookupEnv("FINAM_TOKEN")
	account_id := os.Getenv("FINAM_ACCOUNT_ID")
	_ = account_id

	slog.Info("start")
	// создание клиента
	ctx := context.Background()
	finam.SetLogDebug(true) // проставим признак отладки

	client, err := finam.NewClient(ctx, token)
	if err != nil {
		slog.Error("NewClient", "err", err.Error())
		return
	}

	// получим список всех ордеров по заданному счету
	_ = account_id
	getOrders(ctx, client, account_id)

	// TODO
	// получим информацию по заданному ордеру
	//orderId := "68534984838"
	//getOrder(ctx, client, account_id, orderId)

	// Отмена биржевой заявки
	//orderId := "68534810466" // "665554418102"
	//cancelOrder(ctx, client, account_id, orderId)

	// отмена всех ордеров
	//cancelAllOrders(ctx, client, account_id)

	//getOrders(ctx, client, accountId)
	// еще раз запросим данные по заданному ордеру
	// уже вернется ошибка ("rpc error: code = NotFound desc = Order with id 2033125054207932011 is not found")
	//getOrder(ctx, client, accountId, orderId)

	// пример выставления ордера на покупку\продажу
	//newOrder(ctx, client, account_id)

}

// getOrders получим список всех ордеров по заданному счету
func getOrders(ctx context.Context, client *finam.Client, accountId string) {
	orders, err := client.NewGetOrdersRequest(accountId).Do(ctx)
	if err != nil {
		slog.Error("getOrders", "err", err.Error())
	}
	slog.Info("getOrders", "кол-во ордеров", len(orders.Orders))
	for n, row := range orders.Orders {
		fmt.Println("nrow", n, row.String())
		//slog.Info("OrdersService", "n", n,
		//	"state", row,
		//	"order", row.Order.String())
	}

}

// cancelOrder отмена заявки
func cancelOrder(ctx context.Context, client *finam.Client, accountId, orderId string) {
	orderStatus, err := client.NewCancelOrderRequest(accountId, orderId).Do(ctx)
	if err != nil {
		slog.Error("CancelOrder", "err", err.Error())
		return

	}
	fmt.Println(orderStatus)

}
func cancelAllOrders(ctx context.Context, client *finam.Client, accountId string) {
	err := client.CancelAllOrders(ctx, accountId)
	if err != nil {
		slog.Error("CancelAllOrders", "err", err.Error())
	}
}

// создать новый ордер
func newOrder(ctx context.Context, client *finam.Client, accountId string) {
	symbol := "SRU5@RTSX"
	quantity := 1 // кол-во в штуках

	// покупка по рынку
	//orderStatus, err := client.NewPlaceOrderRequest().AccountId(accountId).Symbol(symbol).Quantity(quantity).Buy().Do(ctx)

	// продажа по рынку
	//orderStatus, err := client.NewPlaceOrderRequest().AccountId(accountId).Symbol(symbol).Quantity(quantity).Sell().Do(ctx)

	// покупка лимитной заявкой
	//orderStatus, err := client.NewPlaceOrderRequest().
	//	AccountId(accountId).Symbol(symbol).Quantity(quantity).BuyLimit().LimitPrice(31992).Do(ctx)

	// продажа лимитной заявкой
	orderStatus, err := client.NewPlaceOrderRequest().
		AccountId(accountId).Symbol(symbol).Quantity(quantity).SellLimit().LimitPrice(31998).Do(ctx)

	if err != nil {
		slog.Error("NewOrder", "err", err.Error())
		return

	}
	fmt.Println(orderStatus)
	//slog.Info("NewOrder", slog.Any("orderStatus", orderStatus))
	//slog.Info("NewOrder", slog.Any("order", orderStatus.Order))
}
