package emailing

import (
	"bytes"

	"github.com/CloudyKit/jet/v6"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	jetset "github.com/neurochar/backend/templates/jet"
)

var tpls = map[string]*jet.Template{}

func NewMessageFromJetTpl(tplFile string, data any) (string, error) {
	tpl, ok := tpls[tplFile]
	if !ok {
		var err error
		tpl, err = jetset.Tpls.GetTemplate(tplFile)
		if err != nil {
			return "", appErrors.ErrInternal.WithWrap(err)
		}
		tpls[tplFile] = tpl
	}

	tplBuffer := new(bytes.Buffer)
	varMap := jet.VarMap{}
	varMap.Set("data", data)

	err := tpl.Execute(tplBuffer, varMap, nil)
	if err != nil {
		return "", appErrors.ErrInternal.WithWrap(err)
	}

	return tplBuffer.String(), nil
}
