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
// список позиций
slog.Info("getPositions", "len(Positions)", len(res.Positions))
for row, pos := range res.Positions {
    slog.Info("AccountsService.GetAccount.Positions",
        "row", row,
        "Symbol", pos.Symbol,
        "Quantity", pos.Quantity.Int(),
        "AveragePrice", pos.AveragePrice.Float64(),
        "CurrentPrice", pos.CurrentPrice.Float64(),
    )

}
```

###  Получение последней котировки по инструменту
```go
symbol := "SBER@MISX" //"ROSN@MISX"  //"SIM5@RTSX"
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
```

### Остальные примеры [тут](/_examples)

## Реализован функционал
### AuthService
1. AuthService.Auth (GetJWT())  
Получение JWT токена из API токена
2. AuthService.TokenDetails (GetTokenDetails)  
Получение информации о токене сессии

### AccountsService
1. AccountsService.GetAccount (AccountRequest)  
Получение информации по конкретному аккаунту
2. AccountsService.Trades (AccountTradesRequest)  
Получение истории по сделкам аккаунта

### AssetsService
0. AssetsService.Clock (GetTimeGetTime)  
Получить текущее время сервера
1. AssetsService.Assets (AssetsRequest)  
Получение списка доступных инструментов, их описание
2. AssetsService.GetAsset (AssetInfoRequest)  
Получение параметров по инструменту
3. AssetsService.GetAssetParams (AssetParamsRequest)  
Получение торговых параметров по инструменту

### MarketDataService
1. MarketDataService.Bars (BarsRequest)  
Получение исторических данных по инструменту (агрегированные свечи)
2. MarketDataService.LastQuote (QuoteRequest)  
Получение последней котировки по инструменту

### OrdersService
1. OrdersService.PlaceOrder (PlaceOrderRequest)  
Выставление биржевой заявки
2. OrdersService.CancelOrder (CancelOrderRequest)        
Отмена биржевой заявки
3. OrdersService.GetOrders (GetOrdersRequest)  
  Получение списка  заявок для аккаунта 


## TODO
* [ ] AssetsService.Schedule
* [ ] MarketDataService.OrderBook
* [ ] OrdersService.GetOrder

