package item

import "github.com/neurochar/backend/pkg/dbhelper"

const (
	Table = "emailing"
)

var TableFields = []string{}

func init() {
	TableFields = dbhelper.ExtractDBFields(&DBModel{})
}
