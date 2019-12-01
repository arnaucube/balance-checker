package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Currency struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Quote  struct {
		USD struct {
			Price float64 `json:"price"`
		}
	}
}
type Content struct {
	Data []Currency
}

type Result struct {
	Symbol  string
	Balance float64
	USD     float64
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "Path to config.toml")
	flag.Parse()

	if err := MustRead(configPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	currencies := getCurrencies()

	results, total := calculateResults(currencies, C.Currencies)
	printResults(results, total)
}

func printResults(results []Result, total float64) {
	fmt.Println("Curr	Balance		USD")
	for _, r := range results {
		fmt.Println(r.Symbol + "	" +
			strconv.FormatFloat(r.Balance, 'f', 6, 64) + "	" +
			strconv.FormatFloat(r.USD, 'f', 6, 64))
	}

	fmt.Println("\ntotal USD:", total)
}

func calculateResults(currencies []Currency, configCurrencies map[string]float64) ([]Result, float64) {
	var results []Result
	var total float64
	total = 0
	for _, c := range currencies {
		if configCurrencies[strings.ToLower(c.Symbol)] != 0 {
			balance := configCurrencies[strings.ToLower(c.Symbol)]
			price := c.Quote.USD.Price
			result := Result{
				Symbol:  c.Symbol,
				Balance: balance,
				USD:     balance * price,
			}
			results = append(results, result)
			total = total + result.USD
		}
	}
	return results, total
}

func getCurrencies() []Currency {
	url := "https://web-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?limit=35&start=1"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var content Content
	err = json.Unmarshal(contents, &content)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return content.Data
}
