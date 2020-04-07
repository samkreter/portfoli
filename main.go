package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/samkreter/portfoli/allocations"
	"github.com/samkreter/portfoli/reader"
)

type Percent float64

// 1. How much money is need to correctly rallocate
// 2. Given a money, best ways to split to go towards allocatoin

func main() {

	filename := "/Users/samkreter/Downloads/portfoli.csv"

	currPositions, err := getCurrentPositions(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Get desired allocation percentages
	desiredAllocation, err := allocations.GetAllocation("Swensen")
	if err != nil {
		log.Fatal(err)
	}

	// Add current asset positions
	for _, asset := range desiredAllocation {
		for _, position := range currPositions {
			if asset.Symbol == position.Symbol {
				asset.CurrValue = position.Current.Value
			}
		}
	}

	// Compute Current Percentages
	desiredAllocation.ComputeCurrPercents()

	// Compute Desired Values based off percentages
	currTotalVal := desiredAllocation.GetCurrTotalVal()
	desiredAllocation.ComputeDesiredValues(currTotalVal)

	// Find the biggest negative off current value
	biggestDiffAsset := desiredAllocation.ComputeGreatestNegativeDiff()
	if biggestDiffAsset == nil {
		log.Println("No Difference")
		return
	}

	// Compute new desired values using new total
	newTotal := biggestDiffAsset.CurrValue / biggestDiffAsset.DesiredPercent
	desiredAllocation.ComputeDesiredValues(newTotal)

	var diff float64
	for _, asset := range desiredAllocation {
		diff = asset.DesiredValue - asset.CurrValue
		fmt.Println(asset.Symbol, "Curr Value: ", asset.CurrValue, "Desired: ", asset.DesiredValue, "Difference: ", diff)
	}

	fmt.Println("Cash required: ", newTotal-currTotalVal)
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
