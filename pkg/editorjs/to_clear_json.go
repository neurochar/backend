package editorjs

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func (d ParsedDocument) ToClearJSON() ([]byte, error) {
	type (
		clearParagraphData struct {
			Text      string             `json:"text"`
			Alignment ParagraphAlignment `json:"alignment"`
		}
		clearHeaderData struct {
			Text  string `json:"text"`
			Level int    `json:"level"`
		}
		clearListItem struct {
			Content string          `json:"content"`
			Meta    map[string]any  `json:"meta"`
			Items   []clearListItem `json:"items"`
		}
		clearListData struct {
			Style ListStyle       `json:"style"`
			Meta  map[string]any  `json:"meta"`
			Items []clearListItem `json:"items"`
		}

		clearImageData struct {
			FileID  *uuid.UUID      `json:"fileId,omitempty"`
			URL     string          `json:"url,omitempty"`
			Type    ImageSourceType `json:"type"`
			Caption string          `json:"caption"`
		}

		clearGalleryFile struct {
			FileID  *uuid.UUID `json:"fileId,omitempty"`
			Caption string     `json:"caption"`
		}
		clearGalleryData struct {
			Files   []clearGalleryFile `json:"files"`
			Caption string             `json:"caption"`
			Style   GalleryStyle       `json:"style"`
		}

		outBlock struct {
			Type string      `json:"type"`
			Data interface{} `json:"data"`
		}
		outDoc struct {
			Time   int64      `json:"time"`
			Blocks []outBlock `json:"blocks"`
		}
	)

	out := outDoc{
		Time:   d.Time.UnixMilli(),
		Blocks: make([]outBlock, 0, len(d.Blocks)),
	}

	var mapListItems func([]ListItem) []clearListItem
	mapListItems = func(items []ListItem) []clearListItem {
		out := make([]clearListItem, 0, len(items))
		for _, it := range items {
			out = append(out, clearListItem{
				Content: it.Content,
				Meta:    it.Meta,
				Items:   mapListItems(it.Items),
			})
		}
		return out
	}

	for _, b := range d.Blocks {
		switch v := b.(type) {
		case ParagraphBlock:
			out.Blocks = append(out.Blocks, outBlock{
				Type: v.Kind(),
				Data: clearParagraphData{
					Text:      v.Data.Text,
					Alignment: v.Data.Alignment,
				},
			})

		case HeaderBlock:
			out.Blocks = append(out.Blocks, outBlock{
				Type: v.Kind(),
				Data: clearHeaderData{
					Text:  v.Data.Text,
					Level: v.Data.Level,
				},
			})

		case ListBlock:
			out.Blocks = append(out.Blocks, outBlock{
				Type: v.Kind(),
				Data: clearListData{
					Style: v.Data.Style,
					Meta:  v.Data.Meta,
					Items: mapListItems(v.Data.Items),
				},
			})

		case ImageBlock:
			if v.Data.File.Type == ImageSrcURL {
				out.Blocks = append(out.Blocks, outBlock{
					Type: v.Kind(),
					Data: clearImageData{
						URL:     v.Data.File.URL,
						Type:    v.Data.File.Type,
						Caption: v.Data.Caption,
					},
				})
			} else {
				out.Blocks = append(out.Blocks, outBlock{
					Type: v.Kind(),
					Data: clearImageData{
						FileID:  v.Data.File.FileID,
						Type:    v.Data.File.Type,
						Caption: v.Data.Caption,
					},
				})
			}

		case GalleryBlock:
			files := make([]clearGalleryFile, 0, len(v.Data.Files))
			for _, f := range v.Data.Files {
				files = append(files, clearGalleryFile{
					FileID:  f.FileID,
					Caption: f.Caption,
				})
			}
			out.Blocks = append(out.Blocks, outBlock{
				Type: v.Kind(),
				Data: clearGalleryData{
					Files:   files,
					Caption: v.Data.Caption,
					Style:   v.Data.Style,
				},
			})

		default:
			return nil, fmt.Errorf("cannot marshal unknown block type %T", b)
		}
	}

	bs, err := json.Marshal(out) // MarshalIndent для красивого вида
	if err != nil {
		return nil, err
	}
	return bs, nil
}
