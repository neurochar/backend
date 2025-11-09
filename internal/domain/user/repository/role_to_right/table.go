package roletoright

import "github.com/neurochar/backend/pkg/dbhelper"

const (
	RoleToRightTable = "role_to_right"
)

var RoleToRightTableFields = []string{}

func init() {
	RoleToRightTableFields = dbhelper.ExtractDBFields(&DBModel{})
}
