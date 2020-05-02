package allocations

import (
	"fmt"
	"log"
	"math"

	"github.com/samkreter/portfoli/asset"
)

// AllocationPlan hold the total plane for all asset allocations
type AllocationPlan struct {
	Name        string
	Allocations []*AssetAllocation
}

// AssetAllocation holds the allocation plan for a single asset
type AssetAllocation struct {
	// NonMutating
	Symbol         string
	DesiredPercent float64

	// Mutatable
	CurrPercent  float64
	CurrValue    float64
	DesiredValue float64
}

// AssetClassPercent shows the percent of an asset class
type AssetClassPercent struct {
	AssetClass    asset.Class
	PercentOfPlan float64
}

// GetAllocation gets an allocation plan by name
func GetAllocation(allocationPlanName string) (AllocationPlan, error) {
	var allocationPlane AllocationPlan

	switch allocationPlanName {
	case "AllWeather":
		allocationPlane = getAllWeatherAllocation()
		if err := allocationPlane.Validate(); err != nil {
			return allocationPlane, err
		}
		return allocationPlane, nil
	case "Swensen":
		allocationPlane = getSwensenAllocation()
		if err := allocationPlane.Validate(); err != nil {
			return allocationPlane, err
		}
		return allocationPlane, nil
	default:
		return allocationPlane, fmt.Errorf("invalid allocation name: %s", allocationPlanName)
	}
}

// Validate ensures the allocation planes is valid
func (plan AllocationPlan) Validate() error {
	totalPercent := 0.0
	for _, aAllocation := range plan.Allocations {
		totalPercent = totalPercent + aAllocation.DesiredPercent
	}

	if math.Round(totalPercent) != 1 {
		return fmt.Errorf("allocation percentage is %d, should be 1", totalPercent)
	}

	return nil
}

// GetCurrTotalVal returns the current total value for the asset plan
func (plan AllocationPlan) GetCurrTotalVal() float64 {
	total := 0.0
	for _, aAllocation := range plan.Allocations {
		total = total + aAllocation.CurrValue
	}

	return total
}

// GetCurrTotalVal returns the current total value for the asset plan
func (plan AllocationPlan) GetDesiredTotalValue() float64 {
	total := 0.0
	for _, aAllocation := range plan.Allocations {
		total = total + aAllocation.DesiredValue
	}

	return total
}

func (plan *AllocationPlan) UpdateDesiredValues() error {
	// Compute Current Percentages
	plan.computeCurrPercents()

	// Compute Desired Values based off percentages
	currTotalVal := plan.GetCurrTotalVal()
	plan.computeDesiredValues(currTotalVal)

	// Find the biggest negative off current value
	biggestDiffAsset := plan.computeGreatestNegativeDiff()
	if biggestDiffAsset.Symbol == "" {
		log.Println("No Difference")
		return nil
	}

	// Compute new desired values using new total
	newTotal := biggestDiffAsset.CurrValue / biggestDiffAsset.DesiredPercent
	plan.computeDesiredValues(newTotal)

	return plan.Validate()
}

// GetAssetClassTotal gets the total asset class percentages for the allocation plan
func (plan AllocationPlan) GetAssetClassTotal() []AssetClassPercent {
	assetClasses := asset.GetAssetClasses()

	classPercents := []AssetClassPercent{}

	for _, class := range assetClasses {
		classPercent := AssetClassPercent{
			AssetClass: class,
		}

		for _, aAllocation := range plan.Allocations {
			a, err := asset.GetAsset(aAllocation.Symbol)
			if err != nil {
				log.Printf("Warning: failed to get asset %q", aAllocation.Symbol)
				continue
			}

			if a.Class == class {
				classPercent.PercentOfPlan += aAllocation.CurrPercent
			}
		}

		classPercents = append(classPercents, classPercent)
	}

	return classPercents
}

// ComputeGreatestNegativeDiff gets the negitive off the current value
func (plan *AllocationPlan) computeGreatestNegativeDiff() AssetAllocation {
	// Find the biggest negative off current value
	biggestDiff := 0.0
	var biggestDiffAsset AssetAllocation

	for idx, aAllocation := range plan.Allocations {
		diff := aAllocation.DesiredValue - aAllocation.CurrValue
		if diff < biggestDiff {
			biggestDiff = diff
			biggestDiffAsset = *plan.Allocations[idx]
		}
	}

	return biggestDiffAsset
}

// ComputeCurrPercents updates the current percents for each asset allocation in the plan
func (plan *AllocationPlan) computeCurrPercents() {
	total := plan.GetCurrTotalVal()

	if math.Round(total) == 0 {
		return
	}

	for _, aAllocation := range plan.Allocations {
		percent := aAllocation.CurrValue / total
		aAllocation.CurrPercent = math.Round(percent*100) / 100
	}
}

// ComputeDesiredValues updates the desired values based on the total and desired percent
func (plan *AllocationPlan) computeDesiredValues(total float64) {
	for _, aAllocation := range plan.Allocations {
		aAllocation.DesiredValue = total * aAllocation.DesiredPercent
	}
}
