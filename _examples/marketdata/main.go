package main

import (
	"context"
	"github.com/Ruvad39/go-finam-rest"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
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
	symbol := "SBER@MISX" //"ROSN@MISX"  //"SIM5@RTSX"

	// Получение последней котировки по инструменту
	getQuote(ctx, client, symbol)
}

// Получение последней котировки по инструменту
func getQuote(ctx context.Context, client *finam.Client, symbol string) {
	quote, err := client.NewQuoteRequest(symbol).Do(ctx)
	if err != nil {
		slog.Error("GetQuote", "err", err.Error())
		return
	}
	slog.Info("GetQuote", "q", quote)
	slog.Info("GetQuote", "symbol", quote.Symbol,
		"Timestamp", quote.Quote.Timestamp.In(finam.TzMoscow),
		"ask", quote.Quote.Ask.Float64(),
		"bid", quote.Quote.Bid.Float64(),
		"last", quote.Quote.Last.Float64(),
		"change", quote.Quote.Change.Float64(),
	)
}
