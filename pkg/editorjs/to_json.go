package editorjs

import (
	"encoding/json"
	"fmt"
)

func (d ParsedDocument) ToJSON() ([]byte, error) {
	type outBlock struct {
		ID   string      `json:"id"`
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}
	type outDoc struct {
		Time    int64      `json:"time"`
		Blocks  []outBlock `json:"blocks"`
		Version string     `json:"version"`
	}

	out := outDoc{
		Time:    d.Time.UnixMilli(),
		Version: d.Version,
		Blocks:  make([]outBlock, 0, len(d.Blocks)),
	}

	for _, b := range d.Blocks {
		switch v := b.(type) {
		case ParagraphBlock:
			out.Blocks = append(out.Blocks, outBlock{
				ID:   v.ID,
				Type: v.Kind(),
				Data: v.Data,
			})
		case HeaderBlock:
			out.Blocks = append(out.Blocks, outBlock{
				ID:   v.ID,
				Type: v.Kind(),
				Data: v.Data,
			})
		case ListBlock:
			out.Blocks = append(out.Blocks, outBlock{
				ID:   v.ID,
				Type: v.Kind(),
				Data: v.Data,
			})
		case ImageBlock:
			out.Blocks = append(out.Blocks, outBlock{
				ID:   v.ID,
				Type: v.Kind(),
				Data: v.Data,
			})
		case GalleryBlock:
			out.Blocks = append(out.Blocks, outBlock{
				ID:   v.ID,
				Type: v.Kind(),
				Data: v.Data,
			})
		default:
			return nil, fmt.Errorf("cannot marshal unknown block type %T", b)
		}
	}

	buf, err := json.Marshal(out)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
