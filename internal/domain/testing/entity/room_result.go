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
	Analyze       *RoomResultAnalyze
}

type RoomResultAnalyzeHiringDecision int

const (
	RoomResultAnalyzeHiringDecisionUnspecified        RoomResultAnalyzeHiringDecision = 0
	RoomResultAnalyzeHiringDecisionHire               RoomResultAnalyzeHiringDecision = 1
	RoomResultAnalyzeHiringDecisionHireWithConditions RoomResultAnalyzeHiringDecision = 2
	RoomResultAnalyzeHiringDecisionDoNotHire          RoomResultAnalyzeHiringDecision = 3
)

type RoomResultAnalyze struct {
	HiringDecision     RoomResultAnalyzeHiringDecision
	ConfidenceScore    float64
	MainRecommendation string
	PersonalityFit     RoomResultAnalyzePersonalityFit
	Risks              []string
	ActionItems        []string
}

type RoomResultAnalyzePersonalityFit struct {
	Score      int
	Summary    string
	KeyMatches []string
	KeyGaps    []string
}
