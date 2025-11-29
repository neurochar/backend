package entity

import "github.com/govalues/decimal"

type RoomResultTechniquesItem struct {
	Result decimal.Decimal
}

type RoomResultTechnique map[uint64]RoomResultTechniquesItem

type RoomResultTraitItem struct {
	TotalResult decimal.Decimal
	Match       decimal.Decimal
	Tip         string
}

type RoomResult struct {
	TotalMatch    decimal.Decimal
	TotalMatchTip string
	Techniques    map[uint64]RoomResultTechnique
	Traits        map[uint64]RoomResultTraitItem
}
