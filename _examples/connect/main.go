package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Ruvad39/go-finam-rest"
	"github.com/joho/godotenv"
)

// предполагаем что есть файл .env в котором записан secret-Token в переменной FINAM_TOKEN
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
	}
}

func main() {
	ctx := context.Background()
	_ = ctx
	token, _ := os.LookupEnv("FINAM_TOKEN")
	slog.Info("start")
	//finam.SetLogDebug(true)

	client, err := finam.NewClient(ctx, token)
	if err != nil {
		slog.Error("NewClient", "err", err.Error())
		return
	}
	// получить текущее время сервера
	t, err := client.GetTime(ctx)
	if err != nil {
		slog.Error("main", "GetTime", err.Error())
		return
	}
	slog.Info("Time", "t", t)
	//return

	//res, err := client.GetJWT()
	// Получение информации о токене сессии
	res, err := client.GetTokenDetails(ctx)
	if err != nil {
		slog.Error("main", "AuthService.TokenDetails", err.Error())
	}

	slog.Info("main", "res", res)
	slog.Info("main", "res.AccountIds", res.AccountIds)
}
