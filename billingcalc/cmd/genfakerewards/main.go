package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
)

const CsvLines = 1000

func main() {
	RunGenFakeRewards()
}

func RunGenFakeRewards() {
	organizations := []string{"Electric Sheep", "Moon Enterprises", "Firetronics"}
	organizationShort := []string{"ES", "ME", "FT"}
	accounts := []string{"Venture Fund", "Capital Fund"}
	operations := []string{"Staking", "Delegation"}
	assets := []string{"CELO", "FLOW", "ROSE", "APT", "SOL", "OSMO", "HASH", "ATOM", "AXL", "EVMOS", "SEI", "SUI", "ETH"}

	rand.Seed(1)

	fileName := "rewards.csv"
	file, err := os.Create(fileName)
	if err != nil {
		return
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"organization", "account", "operation", "asset", "anch_val", "non_anch_val"}
	err = writer.Write(header)
	if err != nil {
		return
	}

	for i := 0; i < CsvLines; i++ {
		var anch_val, non_anch_val float64
		usd := rand.Float64()*100 + 100
		if i%2 == 0 {
			anch_val = usd
		} else {
			non_anch_val = usd
		}

		orgN := rand.Intn(len(organizations))
		data := []string{
			organizations[orgN],
			organizationShort[orgN] + " " + accounts[rand.Intn(len(accounts))],
			operations[rand.Intn(len(operations))],
			assets[rand.Intn(len(assets))],
			fmt.Sprintf("%.2f", anch_val),
			fmt.Sprintf("%.2f", non_anch_val),
		}
		err := writer.Write(data)
		if err != nil {
			return
		}
	}
}
