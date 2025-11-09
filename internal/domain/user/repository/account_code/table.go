package accountcode

import "github.com/neurochar/backend/pkg/dbhelper"

const (
	AccountCodeTable = "account_code"
)

var AccountCodeTableFields = []string{}

func init() {
	AccountCodeTableFields = dbhelper.ExtractDBFields(&DBModel{})
}
