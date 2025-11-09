package file

import "github.com/neurochar/backend/pkg/dbhelper"

const (
	FileTable = "file"
)

var FileTableFields = []string{}

func init() {
	FileTableFields = dbhelper.ExtractDBFields(&DBModel{})
}
