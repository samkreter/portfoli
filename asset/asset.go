package asset

// GetAsset gets an asset by symbol
func GetAsset(symbol string) (Asset, error) {
	asset, ok := knownAssets[symbol]
	if !ok {
		return Asset{}, ErrAssetNotFound
	}

	return asset, nil
}

func GetAssetClasses() []Class {
	return []Class{
		Equity,
		Bond,
		Comodity,
		RealEstate,
	}
}
