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

type command string

func main() {
	filename := flag.String("inputfile", "", "filepath to the fidelity csv (defaults to ~/Downloads/portfoli.csv)")
	assetAllocationName := flag.String("allocation-name", "Swensen", "Name of the asset allocation to use [Swensen, AllWeather]")
	flag.StringVar(assetAllocationName, "a", *assetAllocationName, "Name of the asset allocation to use [Swensen, AllWeather]")
	command := flag.String("c", "printDesired", "the command to use")
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

	// Get desired allocation plan
	allocationPlan, err := allocations.GetAllocation(*assetAllocationName)
	if err != nil {
		log.Fatal(err)
	}

	// Add current asset positions
	for idx, allocationAsset := range allocationPlan.Allocations {
		for _, position := range currPositions {
			if allocationAsset.Symbol == position.Symbol {
				allocationPlan.Allocations[idx].CurrValue = position.Current.Value
			}
		}
	}

	if err := allocationPlan.UpdateDesiredValues(); err != nil {
		log.Fatal(err)
	}

	switch *command {
	case "desired":
		printPlanForReallocation(allocationPlan)
	case "classes":
		printAssetClassPercents(allocationPlan)
	default:
		log.Fatal("Unkown command")
	}

}

func printAssetClassPercents(allocationPlan allocations.AllocationPlan) {
	classPercents := allocationPlan.GetAssetClassTotal()
	for _, classPecent := range classPercents {
		percent := fmt.Sprintf("%d%%", int(classPecent.PercentOfPlan*100))
		fmt.Println(classPecent.AssetClass, ": ", percent)
	}
}

func printPlanForReallocation(allocationPlan allocations.AllocationPlan) {
	var diff float64
	for _, allocationAsset := range allocationPlan.Allocations {
		diff = allocationAsset.DesiredValue - allocationAsset.CurrValue
		fmt.Println(allocationAsset.Symbol, "Curr Value: ", allocationAsset.CurrValue, "Desired: ", allocationAsset.DesiredValue, "Difference: ", diff)
	}

	fmt.Println("Cash required: ", allocationPlan.GetDesiredTotalValue()-allocationPlan.GetCurrTotalVal())
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
