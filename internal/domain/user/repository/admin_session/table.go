package adminsession

import "github.com/neurochar/backend/pkg/dbhelper"

const (
	SessionTable = "auth_admin_session"
)

var SessionTableFields = []string{}

func init() {
	SessionTableFields = dbhelper.ExtractDBFields(&DBModel{})
}
