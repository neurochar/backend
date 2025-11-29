package entity

type Technique interface {
	GetID() uint64
	GetTitle() string
	GetItemsByPersonalityTraits(traitsMap map[uint64]ProfilePersonalityTraitsMapItem) []TechniqueItemData
	CountResult(traitsMap map[uint64]ProfilePersonalityTraitsMapItem, answers map[uint64]any) (RoomResultTechnique, error)
}
