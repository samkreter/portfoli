package allocations

import (
	"fmt"
	"math"
)

type Asset struct {
	Symbol string
	Type string

	DesiredPercent float64

	CurrPercent float64
	CurrValue float64

	DesiredValue float64
}

type AssetAllocation []*Asset

func (a AssetAllocation) Validate() error {
	totalPercent := 0.0
	for _, asset := range a {
		totalPercent = totalPercent + asset.DesiredPercent
	}

	if math.Round(totalPercent) != 1 {
		return fmt.Errorf("allocation percentage is %d, should be 1", totalPercent)
	}

	return nil
}

// Will be nil if no negitive difference
func (a AssetAllocation) ComputeGreatestNegativeDiff() *Asset {
	// Find the biggest negative off current value
	biggestDiff := 0.0
	var biggestDiffAsset *Asset

	for idx, asset := range a {
		diff := asset.DesiredValue - asset.CurrValue
		if diff < biggestDiff {
			biggestDiff = diff
			biggestDiffAsset = a[idx]
		}
	}

	return biggestDiffAsset
}

func (a AssetAllocation) GetCurrTotalVal() float64 {
	total := 0.0
	for _, asset := range a {
		total = total + asset.CurrValue
	}

	return total
}

func (a AssetAllocation) ComputeCurrPercents() {
	total := a.GetCurrTotalVal()

	if math.Round(total) == 0 {
		return
	}

	for _, asset := range a {
		percent := asset.CurrValue / total
		asset.CurrPercent = math.Round(percent * 100) / 100
	}
}

func (a AssetAllocation) ComputeDesiredValues(total float64){
	for _, asset := range a {
		asset.DesiredValue = total * asset.DesiredPercent
	}
}


func GetAllocation(allocationName string) (AssetAllocation, error) {
	var allocation AssetAllocation

	switch allocationName {
	case "AllWeather":
		allocation = getAllWeatherAllocation()
		if err := allocation.Validate(); err !=  nil {
			return allocation, err
		}
		return allocation, nil
	case "Swensen":
		allocation = getSwensenAllocation()
		if err := allocation.Validate(); err !=  nil {
			return allocation, err
		}
		return allocation, nil
	case "Current":
		allocation = getCurrentAllocation()
		if err := allocation.Validate(); err !=  nil {
			return allocation, err
		}
		return allocation, nil
	default:
		return allocation, fmt.Errorf("invalid allocation name: %s", allocationName)
	}
}

func getCurrentAllocation() AssetAllocation {
	return []*Asset{
		{
			Symbol: "VTI",
			Type: "Domestic",
			DesiredPercent: .35,
		},
		{
			Symbol: "VEA",
			Type: "International",
			DesiredPercent: .15,
		},
		{
			Symbol: "VWO",
			Type: "Emerging Markets",
			DesiredPercent: .1,
		},
		{
			Symbol: "TLT",
			Type: "Long Term Treasurey",
			DesiredPercent: .1,
		},
		{
			Symbol: "IEF",
			Type: "Short Term Treasury",
			DesiredPercent: .1,
		},
		{
			Symbol: "DBC",
			Type: "Comodity",
			DesiredPercent: .05,
		},
		{
			Symbol: "GLD",
			Type: "Comodity",
			DesiredPercent: .05,
		},
		{
			Symbol: "VNQ",
			Type: "REITs",
			DesiredPercent: .1,
		},
	}
}

func getAllWeatherAllocation() AssetAllocation {
	return AssetAllocation{
		{
			Symbol: "VTI",
			Type: "Total Stock Market",
			DesiredPercent: .3,
		},
		{
			Symbol: "TLT",
			Type: "20+ Year Treasury",
			DesiredPercent: .4,
		},
		{
			Symbol: "IEF",
			Type: "7-10 Year Treasury",
			DesiredPercent: .15,
		},
		{
			Symbol: "DBC",
			Type: "Commodity Index",
			DesiredPercent: .075,
		},
		{
			Symbol: "GLD",
			Type: "Gold Index",
			DesiredPercent: .075,
		},
	}
}

func getSwensenAllocation() AssetAllocation {
	return AssetAllocation{
		{
			Symbol: "VTI",
			Type: "Domestic",
			DesiredPercent: .3,
		},
		{
			Symbol: "VEA",
			Type: "International",
			DesiredPercent: .15,
		},
		{
			Symbol: "VWO",
			Type: "Emerging Markets",
			DesiredPercent: .05,
		},
		{
			Symbol: "VTIP",
			Type: "US Treasury Inflation Protection Securities",
			DesiredPercent: .15,
		},
		{
			Symbol: "VNQ",
			Type: "REITs",
			DesiredPercent: .2,
		},
		{
			Symbol: "VGIT",
			Type: "U.S. Treasury",
			DesiredPercent: .15,
		},
	}
}