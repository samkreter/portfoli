package allocations

import "fmt"

type Asset struct {
	Symbol string
	Type string
	Percent float64
}

type AssetAllocation []*Asset

func (a AssetAllocation) Validate() error {
	totalPercent := 0.0
	for _, asset := range a {
		totalPercent = totalPercent + asset.Percent
	}

	if totalPercent != 1 {
		return fmt.Errorf("allocation percentage doesn't add to 100%")
	}

	return nil
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
			Percent: .35,
		},
		{
			Symbol: "VEA",
			Type: "International",
			Percent: .15,
		},
		{
			Symbol: "VWO",
			Type: "Emerging Markets",
			Percent: .1,
		},
		{
			Symbol: "TLT",
			Type: "Long Term Treasurey",
			Percent: .1,
		},
		{
			Symbol: "IEF",
			Type: "Short Term Treasury",
			Percent: .1,
		},
		{
			Symbol: "DBC",
			Type: "Comodity",
			Percent: .05,
		},
		{
			Symbol: "GLD",
			Type: "Comodity",
			Percent: .05,
		},
		{
			Symbol: "VNQ",
			Type: "REITs",
			Percent: .1,
		},
	}
}

func getAllWeatherAllocation() AssetAllocation {
	return AssetAllocation{
		{
			Symbol: "VTI",
			Type: "Total Stock Market",
			Percent: .3,
		},
		{
			Symbol: "TLT",
			Type: "20+ Year Treasury",
			Percent: .4,
		},
		{
			Symbol: "IEF",
			Type: "7-10 Year Treasury",
			Percent: .15,
		},
		{
			Symbol: "DBC",
			Type: "Commodity Index",
			Percent: .075,
		},
		{
			Symbol: "GLD",
			Type: "Gold Index",
			Percent: .075,
		},
	}
}

func getSwensenAllocation() AssetAllocation {
	return AssetAllocation{
		{
			Symbol: "VTI",
			Type: "Domestic",
			Percent: .30,
		},
		{
			Symbol: "VEA",
			Type: "International",
			Percent: .15,
		},
		{
			Symbol: "VWO",
			Type: "Emerging Markets",
			Percent: .05,
		},
		{
			Symbol: "VTIP",
			Type: "US Treasury Inflation Protection Securities",
			Percent: .15,
		},
		{
			Symbol: "VNQ",
			Type: "REITs",
			Percent: .2,
		},
		{
			Symbol: "VGIT",
			Type: "U.S. Treasury",
			Percent: .15,
		},
	}
}