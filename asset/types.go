package asset

import "errors"

// Class the asset class
type Class string

const (
	Equity     = Class("Equity")
	Bond       = Class("Bond")
	Comodity   = Class("Commodities")
	RealEstate = Class("RealEstate")
)

// SubClass an asset subclass
type SubClass string

const (
	// Equity
	Domestic        = SubClass("Domestic")
	International   = SubClass("International")
	EmergingMarkets = SubClass("Emerging Markets")

	// Bonds
	LongTermTreasury             = SubClass("Long Term Treasurey")
	MediumTermTreasury           = SubClass("Medium Term Treasury")
	InflationProtectedSecurities = SubClass("US Treasury Inflation Protection Securities")

	// Comodity
	Gold  = SubClass("Gold")
	Index = SubClass("Index")

	// Realestate
	Reits = SubClass("REITs")
)

// Asset holds information for a specific asset
type Asset struct {
	Symbol   string
	Class    Class
	SubClass SubClass
}

var (
	knownAssets = map[string]Asset{
		"VTI": {
			Symbol:   "VTI",
			Class:    Equity,
			SubClass: Domestic,
		},
		"VEA": {
			Symbol:   "VEA",
			Class:    Equity,
			SubClass: International,
		},
		"VWO": {
			Symbol:   "VWO",
			Class:    Equity,
			SubClass: EmergingMarkets,
		},
		"TLT": {
			Symbol:   "TLT",
			Class:    Bond,
			SubClass: LongTermTreasury,
		},
		"IEF": {
			Symbol:   "IEF",
			Class:    Bond,
			SubClass: MediumTermTreasury,
		},
		"DBC": {
			Symbol:   "DBC",
			Class:    Comodity,
			SubClass: Index,
		},
		"GLD": {
			Symbol:   "GLD",
			Class:    Comodity,
			SubClass: Gold,
		},
		"VNQ": {
			Symbol:   "VNQ",
			Class:    RealEstate,
			SubClass: Reits,
		},
		"VTIP": {
			Symbol:   "VTIP",
			Class:    Bond,
			SubClass: InflationProtectedSecurities,
		},
		"VGIT": {
			Symbol:   "VGIT",
			Class:    Bond,
			SubClass: MediumTermTreasury,
		},
	}
)

var (
	ErrAssetNotFound = errors.New("Asset Not Found")
)
