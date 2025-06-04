# go-finam-rest

**rest-клиент на Go для работы с API Финама**  
[tradeapi.finam](https://tradeapi.finam.ru/docs/about/)

## Установка

```bash
go get github.com/Ruvad39/go-finam-rest
```

## Примеры

### Пример создание клиента. Получение данных по токену
```go
ctx := context.Background()
token, _ := "FINAM_TOKEN"
client, err := finam.NewClient(ctx, token)
if err != nil {
    slog.Error("NewClient", "err", err.Error())
return
}

// Получение информации о токене сессии
res, err := client.GetTokenDetails(ctx)
if err != nil {
slog.Error("main", "AuthService.TokenDetails", err.Error())
}
slog.Info("main", "res", res)
slog.Info("main", "res.AccountIds", res.AccountIds)

```

###  Получить информацию по торговому счету
```go
accountId := "номер счета"
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
)
```


### Примеры смотрите [тут](/_examples)


## TODO
* [x] AuthService
* [x] OrdersService
* [ ] AccountsService
* [ ] AssetsService
* [ ] MarketDataService
