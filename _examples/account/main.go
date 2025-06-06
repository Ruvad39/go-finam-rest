package main

import (
	"context"
	"fmt"
	"github.com/Ruvad39/go-finam-rest"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"time"
)

func main() {
	// предполагаем что есть файл .env в котором записан secret-Token в переменной FINAM_TOKEN
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
	}
	ctx := context.Background()
	token, _ := os.LookupEnv("FINAM_TOKEN")
	slog.Info("start")
	finam.SetLogDebug(true)

	client, err := finam.NewClient(ctx, token)
	if err != nil {
		slog.Error("NewClient", "err", err.Error())
		return
	}

	// account_Id, _ := os.LookupEnv("FINAM_ACCOUNT_ID")
	// getAccount(ctx, client, account_Id)
	// getPositions(ctx, client, account_Id)
	// getTrades(ctx, client, account_Id)
	// return

	// Получение информации о токене сессии. Возьмем список счетов
	res, err := client.GetTokenDetails(ctx)
	if err != nil {
		slog.Error("main", "AuthService.TokenDetails", err.Error())
	}
	for row, accountId := range res.AccountIds {
		// Получение информации по конкретному аккаунту
		slog.Info("TokenDetails.AccountIds", "row", row, "accoiuntId", accountId)
		// получим информацию по конкретному счету
		//getAccount(ctx, client, accountId)
		//getPositions(ctx, client, accountId)
		getTrades(ctx, client, accountId)
		//getTransactions(ctx, client, accountId)
	}
}

// getAccount получим информацию по конкретному счету
func getAccount(ctx context.Context, client *finam.Client, accountId string) {
	res, err := client.NewAccountRequest(accountId).Do(ctx)
	if err != nil {
		slog.Error("AccountsService.GetAccount", "GetAccount", err.Error())
		return
	}
	slog.Info("AccountsService.GetAccount",
		"AccountId", res.AccountId,
		"Type", res.Type,
		"Status", res.Status,
		"Equity", fmt.Sprintf("%.2f", res.Equity.Float64()),
		"UnrealizedProfit", fmt.Sprintf("%.2f", res.UnrealizedProfit.Float64()),
		"Cash", res.Cash[0], // будет ошибка, если нет денег
		//"Cash", res.Cash[0].String(),
	)
}

func getPositions(ctx context.Context, client *finam.Client, accountId string) {
	res, err := client.NewAccountRequest(accountId).Do(ctx)
	if err != nil {
		slog.Error("AccountsService.GetAccount", "GetAccount", err.Error())
		return
	}
	slog.Info("getPositions", "len(Positions)", len(res.Positions))
	// список позиций
	for row, pos := range res.Positions {
		slog.Info("AccountsService.GetAccount.Positions",
			"row", row,
			"Symbol", pos.Symbol,
			"Quantity", pos.Quantity.Int(),
			"AveragePrice", pos.AveragePrice.Float64(),
			"CurrentPrice", pos.CurrentPrice.Float64(),
		)

	}
}

func getTrades(ctx context.Context, client *finam.Client, accountId string) {
	// запросим все сделки за последние 24 часа
	var limit int32 = 0
	start_time := time.Now().Add(-24 * time.Hour) //  24 часа назад
	end_time := time.Now()
	res, err := client.NewAccountTradesRequest(accountId).Limit(limit).StartTime(start_time).EndTime(end_time).Do(ctx)
	if err != nil {
		slog.Error("AccountsService.GetTrades", "GetTrades", err.Error())
		return
	}

	slog.Info("getTrades", "len(Trades)", len(res.Trades))
	// список сделок
	for row, t := range res.Trades {
		slog.Info("AccountsService.GetAccount.Trades",
			"row", row,
			"Trade", t,
		)

	}
}
