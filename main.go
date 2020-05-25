package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/samkreter/portfoli/allocations"
	"github.com/samkreter/portfoli/pkg/fidelity"
)

type command string

func main() {
	filename := flag.String("inputfile", "", "filepath to the fidelity csv (defaults to ~/Downloads/portfoli.csv)")

	assetAllocationName := flag.String("allocation-name", "Swensen", "Name of the asset allocation to use [Swensen, AllWeather]")
	flag.StringVar(assetAllocationName, "a", *assetAllocationName, "Name of the asset allocation to use [Swensen, AllWeather]")

	command := flag.String("c", "desired", "the command to use")
	flag.Parse()

	currPositions, err := fidelity.GetCurrentPositions(*filename)
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
