package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"

	"github.com/samkreter/portfoli/allocations"
	"github.com/samkreter/portfoli/reader"
)

// 1. How much money is need to correctly rallocate
// 2. Given a money, best ways to split to go towards allocatoin

func main() {
	filename := flag.String("inputfile", "", "filepath to the fidelity csv (defaults to ~/Downloads/portfoli.csv)")
	assetAllocationName := flag.String("allocation-name", "Swensen", "Name of the asset allocation to use [Swensen, AllWeather]")
	flag.Parse()

	if *filename == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		*filename = usr.HomeDir + "/Downloads/portfoli.csv"
	}

	currPositions, err := getCurrentPositions(*filename)
	if err != nil {
		log.Fatal(err)
	}

	// Get desired allocation percentages
	desiredAllocation, err := allocations.GetAllocation(*assetAllocationName)
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

	rows, err := readCSV(f)
	if err != nil {
		return nil, err
	}

	currentPositions := []*reader.FidelityRow{}
	for _, row := range rows {
		fRow := reader.ParseRow(row)
		if fRow != nil {
			currentPositions = append(currentPositions, fRow)
		}
	}

	return currentPositions, nil
}

func readCSV(f *os.File) ([][]string, error) {
	rows := [][]string{}
	log.Println("here")
	csvReader := csv.NewReader(f)
	csvReader.LazyQuotes = true

	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				return rows, nil
			}
			// ignore invalid filed count since the fidelity csv's have junk at the end
			if err, ok := err.(*csv.ParseError); ok && err.Err == csv.ErrFieldCount {
				continue
			}

			return nil, err
		}

		rows = append(rows, row)
	}
}
