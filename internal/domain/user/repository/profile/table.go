package profile

import "github.com/neurochar/backend/pkg/dbhelper"

const (
	ProfileTable = "profile"
)

var ProfileTableFields = []string{}

func init() {
	ProfileTableFields = dbhelper.ExtractDBFields(&DBModel{})
}
