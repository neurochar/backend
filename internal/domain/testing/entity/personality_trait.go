package entity

type PersonalityTraitType int

const (
	PersonalityTraitTypeUnspecified PersonalityTraitType = 0
	PersonalityTraitTypeBipolar     PersonalityTraitType = 1
)

var PersonalityTraitTypeMap = map[PersonalityTraitType]string{
	PersonalityTraitTypeBipolar: "bipolar",
}

func (t PersonalityTraitType) String() string {
	return PersonalityTraitTypeMap[t]
}

type PersonalityTrait interface {
	GetID() uint64
	GetType() PersonalityTraitType
	GetName() string
	GetDescription() string
	GetLeftStateName() string
	GetRightStateName() string
}

type PersonalityTraitBipolar struct {
	ID             uint64
	Name           string
	Description    string
	LeftStateName  string
	RightStateName string
}

func (t *PersonalityTraitBipolar) GetID() uint64 {
	return t.ID
}

func (t *PersonalityTraitBipolar) GetType() PersonalityTraitType {
	return PersonalityTraitTypeBipolar
}

func (t *PersonalityTraitBipolar) GetName() string {
	return t.Name
}

func (t *PersonalityTraitBipolar) GetDescription() string {
	return t.Description
}

func (t *PersonalityTraitBipolar) GetLeftStateName() string {
	return t.LeftStateName
}

func (t *PersonalityTraitBipolar) GetRightStateName() string {
	return t.RightStateName
}
