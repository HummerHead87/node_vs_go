package main

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/adshao/go-binance"
)

type BinanceSymbol struct {
	pair  string
	price float64
}

func main() {
	coinsType := [...]string{"BTC",
		"BNB",
		"ETH",
		"PAX",
		"USDC",
		"USDT",
		"TUSD",
		"USDS",
		"XRP"}

	prices, err := getBinanceExchanges()
	if err != nil {
		fmt.Println(err)
	}

	start := time.Now()

	result := make([]BinanceSymbol, len(prices))

	for i, p := range prices {
		if price, err := strconv.ParseFloat(p.Price, 64); err == nil {
			for _, coin := range coinsType {
				re := regexp.MustCompile(fmt.Sprintf("(.*)(%s$)", coin))
				if match := re.FindAllStringSubmatch(p.Symbol, 1); match != nil {
					pair := fmt.Sprintf("%s_%s", match[0][1], match[0][2])

					result[i] = BinanceSymbol{
						pair:  pair,
						price: price,
					}

					break
				}
			}
		}
	}

	sort.Slice(result, func(i, j int) bool { return result[i].pair < result[j].pair })
	fmt.Println(time.Since(start))
	// fmt.Println(result)
}

func getBinanceExchanges() ([]*binance.SymbolPrice, error) {
	var (
		apiKey    = ""
		secretKey = ""
	)

	client := binance.NewClient(apiKey, secretKey)

	prices, err := client.NewListPricesService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	// for _, p := range prices {
	// 	fmt.Println(p)
	// }

	return prices, nil
}
