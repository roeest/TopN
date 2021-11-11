package main

import (
	"sort"
	"flag"
	"fmt"
	"log"
	"math"

	"github.com/roeest/vanguard"
)

type sectorDetails struct {
	description string
}

var sectorEtfs = map[string]string{
	"VOX": "Communication Services",
	"VCR": "Consumer Discretionary",
	"VDC": "Consumer Staples",
	"VDE": "Energy",
	"VFH": "Financials",
	"VHT": "Health Care",
	"VIS": "Industrials",
	"VGT": "Information Technology",
	"VAW": "Materials",
	"VNQ": "Real Estate",
	"VPU": "Utilities",
}

const snp500ETF = "VOO"

func main() {
	amount := flag.Int("amount", 10000, "the maount of money to invest")
	n := flag.Int("n", 3, "the top stocks in each sector")
	flag.Parse()
	cli := vanguard.NewClient()

	voo, err := cli.GetEtf(snp500ETF)
	if err != nil {
		log.Fatalln("failed getting S&P500 etf", err)
	}
	di, err := voo.GetDiversificationInfo()
	if err != nil {
		log.Fatalln("failed getting diversification info for VOO", err)
	}

	count := 1
	toInvest := 0.0
	for symbol, description := range sectorEtfs {
		sectorToInvest := 0
		etf, err := cli.GetEtf(symbol)
		if err != nil {
			log.Fatalln("failed getting etf", err)
		}

		holdings, err := etf.GetHoldings(vanguard.Stock)
		if err != nil {
			log.Fatalln("failed getting holdings", err)
		}
		topHoldings, marketSize := getHoldings(holdings, *n)

		weight := di.Sectors[description].BenchmarkWeight
		for _, h := range topHoldings {
			holding := math.Round(weight * (h.value / marketSize) * float64(*amount))
			toInvest += holding
			fmt.Printf("%d. (%v) Buy %v %f$\n", count, description, h.symbol, holding)
			count++
			sectorToInvest += int(holding)
		}
		defer fmt.Printf("Total invested in %s: %d (weight: %f%% - expected %d) \n", description, sectorToInvest, 100.0*weight, int64(weight*float64(*amount)))
	}
	defer fmt.Println("Total:", toInvest)
}

type holding struct {
	symbol string
	value  float64
}

type holdings []holding

func (h holdings) Len() int           { return len(h) }
func (h holdings) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h holdings) Less(i, j int) bool { return h[i].value > h[j].value }

func getHoldings(holdings []vanguard.Holding, n int) (holdings, float64) {
	holdingMap := make(map[string]float64, n)
	marketSize := 0
	for _, h := range holdings {
		if h.Symbol == "" {
			continue
		}
		_, ok := holdingMap[h.Symbol]
		if len(holdingMap) == n {
			if !ok {
				return toOrderedSlice(holdingMap), float64(marketSize)
			}
		}
		marketSize += h.MarketValue
		if ok {
			holdingMap[h.Symbol] += float64(h.MarketValue)
			continue
		}

		holdingMap[h.Symbol] = float64(h.MarketValue)
	}
	return nil, 0
}

func toOrderedSlice(h map[string]float64) holdings {
	result := make(holdings, 0, len(h))
	for k, v := range h {
		result = append(result, holding{symbol: k, value: v})
	}
	sort.Sort(result)
	return result
}