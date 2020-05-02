package allocations

func getAllWeatherAllocation() AllocationPlan {
	return AllocationPlan{
		Name: "AllWeather",
		Allocations: []*AssetAllocation{
			{
				Symbol:         "VTI",
				DesiredPercent: .3,
			},
			{
				Symbol:         "TLT",
				DesiredPercent: .4,
			},
			{
				Symbol:         "IEF",
				DesiredPercent: .15,
			},
			{
				Symbol:         "DBC",
				DesiredPercent: .075,
			},
			{
				Symbol:         "GLD",
				DesiredPercent: .075,
			},
		},
	}
}

func getSwensenAllocation() AllocationPlan {
	return AllocationPlan{
		Name: "Swensen",
		Allocations: []*AssetAllocation{
			{
				Symbol:         "VTI",
				DesiredPercent: .3,
			},
			{
				Symbol:         "VEA",
				DesiredPercent: .15,
			},
			{
				Symbol:         "VWO",
				DesiredPercent: .05,
			},
			{
				Symbol:         "VTIP",
				DesiredPercent: .15,
			},
			{
				Symbol:         "VNQ",
				DesiredPercent: .2,
			},
			{
				Symbol:         "VGIT",
				DesiredPercent: .15,
			},
		},
	}
}
