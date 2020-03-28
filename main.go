package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/samkreter/portfoli/reader"
	"github.com/samkreter/portfoli/allocations"
)

type Percent float64


// 1. How much money is need to correctly rallocate
// 2. Given a money, best ways to split to go towards allocatoin

func main() {

	filename := "/Users/samkreter/Downloads/Portfolio_Position_Mar-27-2020.csv"

	currPositions, err := getCurrentPositions(filename)
	if err != nil {
		log.Fatal(err)
	}

	desiredAllocation, err := allocations.GetAllocation("Swensen")
	if err != nil {
		log.Fatal(err)
	}

	for _, asset := range desiredAllocation {
		for _, position := range currPositions {
			if asset.Symbol == position.Symbol {
				asset.CurrValue = position.Current.Value
			}
		}
	}

	// Setup current percentages
	desiredAllocation.ComputeCurrPercents()

	// Setup Desired Values
	currTotalVal := desiredAllocation.GetCurrTotalVal()
	desiredAllocation.ComputeDesiredValues(currTotalVal)

	// Find the biggest negitive to set as new base
	biggestDiff := 0.0
	currentVal := 0.0
	desiredPecent := 0.0
	for _, asset := range desiredAllocation {
		diff := asset.DesiredValue - asset.CurrValue
		if diff < biggestDiff {
			biggestDiff = diff
			currentVal = asset.CurrValue
			desiredPecent = asset.DesiredPercent
		}
	}

	if biggestDiff == 0.0 {
		log.Println("No difference")
		return
	}

	newTotal := currentVal / desiredPecent

	desiredAllocation.ComputeDesiredValues(newTotal)

	for _, asset := range desiredAllocation {
		fmt.Println(asset.Symbol, "Curr Value: ", asset.CurrValue, "Desired: ", asset.DesiredValue)
	}

	fmt.Println("Cash required: ", newTotal - currTotalVal)
}

func getCurrentPositions(filename string) ([]*reader.FidelityRow, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	currentPositions := []*reader.FidelityRow{}
	for _, line := range lines {
		row := reader.ParseRow(line)
		if row != nil {
			currentPositions = append(currentPositions, row)
		}
	}

	return currentPositions, nil
}
