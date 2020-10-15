package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	"github.com/roeest/vanguard"
)

type sectorDetails struct {
	weight      float64
	description string
}

var sectors = map[string]sectorDetails{
	"VOX": {weight: 0.1, description: "Communication Services"},
	"VCR": {weight: 0.102, description: "Consumer Discretionary"},
	"VDC": {weight: 0.067, description: "Consumer Staples"},
	"VDE": {weight: 0.06, description: "Energy"},
	"VFH": {weight: 0.137, description: "Financials"},
	"VHT": {weight: 0.149, description: "Health Care"},
	"VIS": {weight: 0.097, description: "Industrials"},
	"VGT": {weight: 0.208, description: "Information Technology"},
	"VAW": {weight: 0.025, description: "Materials"},
	"VNQ": {weight: 0.027, description: "Real Estate"},
	"VPU": {weight: 0.028, description: "Utilities"},
}

func main() {
	amount := flag.Int("amount", 10000, "the maount of money to invest")
	n := flag.Int("n", 3, "the top stocks in each sector")
	flag.Parse()
	cli := vanguard.NewClient()
	count := 1
	toInvest := 0.0
	for symbol, details := range sectors {
		etf, err := cli.GetEtf(symbol)
		if err != nil {
			log.Fatalln("failed getting etf", err)
		}

		holdings, err := etf.GetHoldings(vanguard.Stock)
		if err != nil {
			log.Fatalln("failed getting holdings", err)
		}
		topHoldings, marketSize := getHoldings(holdings, *n)

		for symbol, value := range topHoldings {
			holding := math.Round(details.weight * (value / marketSize) * float64(*amount))
			toInvest += holding
			fmt.Printf("%d. (%v) Buy %v %f$\n", count, details.description, symbol, holding)
			count++
		}
	}
	fmt.Println("Total:", toInvest)
}

func getHoldings(holdings []vanguard.Holding, n int) (map[string]float64, float64) {
	result := make(map[string]float64, n)
	marketSize := 0
	for _, h := range holdings {
		if h.Symbol == "" {
			continue
		}
		_, ok := result[h.Symbol]
		if len(result) == n {
			if !ok {
				return result, float64(marketSize)
			}
		}
		marketSize += h.MarketValue
		if ok {
			result[h.Symbol] += float64(h.MarketValue)
			continue
		}

		result[h.Symbol] = float64(h.MarketValue)
	}
	return nil, 0
}
