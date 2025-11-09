package entity

type RightType int

const (
	RightTypeBool RightType = 1
	RightTypeInt  RightType = 2
)

func (r RightType) String() string {
	return map[RightType]string{
		RightTypeBool: "bool",
		RightTypeInt:  "int",
	}[r]
}

type Right struct {
	ID                uint64
	Key               string
	Type              RightType
	DefaultValue      int
	DefaultSuperValue int
}
