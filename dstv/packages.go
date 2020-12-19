package main

import "sort"

const (
	ACCESS      = "access"
	PREMIUM     = "premium"
	COMPACT     = "compact"
	FAMILY      = "family"
	COMPACTPLUS = "compact plus"
)

var currency string = "ssp"

const rate float64 = 550.0

type Package struct {
	Name     string  `json:"package"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	SerialNo int     `json:"serial_no"`
}

func getPackages() []Package {
	access := Package{}
	premium := Package{}
	family := Package{}
	compactdstv := Package{}
	compactplus := Package{}

	//access
	access.Name = ACCESS
	access.Price = (12.0 * rate)
	access.Currency = currency
	access.SerialNo = 1

	//premium
	premium.Name = PREMIUM
	premium.Price = (82.0 * rate)
	premium.Currency = currency
	premium.SerialNo = 5

	//FAMILY
	family.Name = FAMILY
	family.Price = (19.0 * rate)
	family.Currency = currency
	family.SerialNo = 2

	//compactplus
	compactplus.Name = COMPACTPLUS
	compactplus.Price = (53.0 * 320)
	compactplus.Currency = currency
	compactplus.SerialNo = 4

	//compact
	compactdstv.Name = COMPACT
	compactdstv.Price = (25.0 * 320)
	compactdstv.Currency = currency
	compactdstv.SerialNo = 3
	packs := []Package{access, family, compactdstv, compactplus, premium}
	sort.Slice(packs, func(i int, j int) bool {
		return packs[i].SerialNo < packs[j].SerialNo
	})
	return packs
}
