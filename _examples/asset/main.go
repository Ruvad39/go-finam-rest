package main

import (
	"context"
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

	// Получение списка доступных инструментов, их описание
	//getAssets(ctx, client)

	//
	//getAssetInfo(ctx, client, "SBER@MISX", account_id)
	symbol := "FEES@MISX" //"EDU5@RTSX" // "SRM5@RTSX" // "SBER@MISX" "FEES@MISX"
	getAssetInfo(ctx, client, symbol, account_id)
	//slog.Info("account_id",account_id)
	//getAssetParams(ctx, client, symbol, account_id)

	// TODO Получение расписания торгов для инструмента
	//getSchedule(ctx, client, "SBER@MISX")

}

// Получение списка доступных инструментов, их описание
func getAssets(ctx context.Context, client *finam.Client) {

	assets, err := client.NewAssetsRequest().Do(ctx)
	if err != nil {
		slog.Error("AssetsRequest", "err", err.Error())
	}
	slog.Info("AssetsRequest", "assets.len", len(assets.Assets))

	for n, sec := range assets.Assets {
		//if sec.Type == "FUTURES" && sec.Mic == "RTSX" && 1 == 2 {
		if 1 == 1 {
			slog.Info("assets",
				"row", n,
				"id", sec.Id,
				"Symbol", sec.Symbol,
				"Name", sec.Name,
				"ticker", sec.Ticker,
				"Type", sec.Type,
				"Mic", sec.Mic,
			)
		}
	}
}

// Получение списка доступных инструментов, их описание
func getAssetInfo(ctx context.Context, client *finam.Client, symbol, accountId string) {

	info, err := client.NewAssetInfoRequest(symbol, accountId).Do(ctx)
	if err != nil {
		slog.Error("AssetsRequest", "err", err.Error())
	}
	slog.Info("AssetsInfoRequest", "info", info)
}

// Получение списка доступных инструментов, их описание
func getAssetParams(ctx context.Context, client *finam.Client, symbol, accountId string) {

	err := client.NewAssetParamsRequest(symbol, accountId).Do(ctx)
	if err != nil {
		slog.Error("AssetsRequest", "err", err.Error())
	}
}
