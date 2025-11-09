package role

import "github.com/neurochar/backend/pkg/dbhelper"

const (
	RoleTable = "role"
)

var RoleTableFields = []string{}

func init() {
	RoleTableFields = dbhelper.ExtractDBFields(&DBModel{})
}
