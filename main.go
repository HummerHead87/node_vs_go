package main

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/adshao/go-binance"
)

var coinsType = [...]string{"BTC",
	"BNB",
	"ETH",
	"PAX",
	"USDC",
	"USDT",
	"TUSD",
	"USDS",
	"XRP"}

type BinanceSymbol struct {
	pair  string
	price float64
}

type BencmarkResult struct {
	min time.Duration
	max time.Duration
	avg time.Duration
}

func (b BencmarkResult) String() string {
	return fmt.Sprintf("min: %s, max: %s, avg: %s", b.min, b.max, b.avg)
}

func main() {
	prices, err := getBinanceExchanges()
	if err != nil {
		fmt.Println(err)
	}

	// start := time.Now()

	// result := formatInput(prices)
	// fmt.Println(time.Since(start))
	// fmt.Println(result)
	fmt.Println(benchMark(1000, formatInput, prices))
}

func benchMark(iterations int, f func([]*binance.SymbolPrice) []BinanceSymbol, args []*binance.SymbolPrice) BencmarkResult {
	times := make([]time.Duration, iterations)

	for i := 0; i < iterations; i++ {
		start := time.Now()
		result := f(args)
		if result != nil {
			times[i] = time.Since(start)
		}
	}

	return BencmarkResult{
		min: minDuration(times),
		max: maxDuration(times),
		avg: avgDuration(times),
	}
}

func minDuration(values []time.Duration) time.Duration {
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}

	return min
}

func maxDuration(values []time.Duration) time.Duration {
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}

	return max
}

func avgDuration(values []time.Duration) time.Duration {
	var summ time.Duration

	for _, v := range values {
		summ += v
	}

	avg := int(summ/time.Microsecond) / len(values)
	// return summ / (len(values) * time.Microsecond)
	return time.Duration(avg) * time.Microsecond
}

func formatInput(prices []*binance.SymbolPrice) []BinanceSymbol {
	result := make([]BinanceSymbol, len(prices))

	for i, p := range prices {
		if price, err := strconv.ParseFloat(p.Price, 64); err == nil {
			result[i] = BinanceSymbol{
				pair:  _splitSymbolStr2(p.Symbol),
				price: price,
			}
		}
	}

	sort.Slice(result, func(i, j int) bool { return result[i].pair < result[j].pair })

	return result
}

func _splitSymbolStr(symbol string) string {
	var result string

	for _, coin := range coinsType {
		re := regexp.MustCompile(fmt.Sprintf("(.*)(%s$)", coin))
		if match := re.FindAllStringSubmatch(symbol, 1); match != nil {
			result = fmt.Sprintf("%s_%s", match[0][1], match[0][2])

			break
		}
	}

	return result
}

func _splitSymbolStr2(symbol string) string {
	var result string

	for _, coin := range coinsType {
		if strings.HasSuffix(symbol, coin) {
			i := utf8.RuneCountInString(symbol) - utf8.RuneCountInString(coin)
			symbolSlice := []rune(symbol)
			result = string(append(symbolSlice[:i], append([]rune{'_'}, symbolSlice[i:]...)...))
			// result = fmt.Sprintf("%s_%s", string(symbolSlice[:i]), string(symbolSlice[i:]))
			break
		}
	}

	return result
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
