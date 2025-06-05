package main

import (
	"context"
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
	symbol := "SBER@MISX" //"ROSN@MISX"  //"SIM5@RTSX"

	// Получение последней котировки по инструменту
	//getQuote(ctx, client, symbol)

	// получение списка свечей
	getBars(ctx, client, symbol)
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

// получение списка свечей
func getBars(ctx context.Context, client *finam.Client, symbol string) {
	//symbol := "SBER@MISX" //"SIM5@RTSX" MISX
	// получение списка свечей

	start_time, _ := time.Parse("2006-01-02", "2025-05-01")
	end_time := time.Now()
	bars, err := client.NewBarsRequest().Symbol(symbol).Timeframe(finam.TimeframeD1).StartTime(start_time).EndTime(end_time).Do(ctx)

	if err != nil {
		slog.Error("MarketDataService.Bars", "err", err.Error())
		return
	}
	slog.Info("MarketDataService.Bars", "Bars.len", len(bars.Bars))
	for row, bar := range bars.Bars {
		slog.Info("Bar", "row", row,
			"Timestamp", bar.Timestamp.In(finam.TzMoscow),
			"Open", bar.Open.Float64(),
			"High", bar.High.Float64(),
			"Low", bar.Low.Float64(),
			"Close", bar.Close.Float64(),
			"Volume", bar.Volume.Int(),
			"TF", bars.Timeframe,
		)
	}
}
