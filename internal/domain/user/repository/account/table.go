package account

import "github.com/neurochar/backend/pkg/dbhelper"

const (
	AccountTable = "account"
)

var AccountTableFields = []string{}

func init() {
	AccountTableFields = dbhelper.ExtractDBFields(&DBModel{})
}
