package editorjs

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/google/uuid"
)

func ParseEditorJS(input []byte) (ParsedDocument, []uuid.UUID, error) {
	var rd rawDocument
	if err := json.Unmarshal(input, &rd); err != nil {
		return ParsedDocument{}, nil, fmt.Errorf("invalid JSON: %w", err)
	}

	var errs MultiError
	// if rd.Version == "" {
	// 	errs.Add("version", "must be non-empty")
	// }
	// if len(rd.Blocks) == 0 {
	// 	errs.Add("blocks", "must contain at least 1 block")
	// }

	if err := errs.OrNil(); err != nil {
		return ParsedDocument{}, nil, err
	}

	doc := ParsedDocument{
		Time:    rd.Time.Time(),
		Version: rd.Version,
		Blocks:  make([]AnyBlock, 0, len(rd.Blocks)),
	}

	files := make([]uuid.UUID, 0)
	addFile := func(id *uuid.UUID) {
		if id == nil || *id == uuid.Nil {
			return
		}
		if !slices.Contains(files, *id) {
			files = append(files, *id)
		}
	}

	for i, rb := range rd.Blocks {
		path := fmt.Sprintf("blocks[%d]", i)

		if rb.ID == "" {
			errs.Add(path+".id", "required")
			continue
		}
		if rb.Type == "" {
			errs.Add(path+".type", "required")
			continue
		}
		if len(rb.Data) == 0 {
			errs.Add(path+".data", "required")
			continue
		}

		switch rb.Type {
		case "paragraph":
			var block ParagraphBlock
			block.ID = rb.ID
			if err := json.Unmarshal(rb.Data, &block.Data); err != nil {
				errs.Add(path+".data", "invalid paragraph payload: "+err.Error())
				continue
			}
			if e := validateParagraph(&block.Data, path+".data"); e != nil {
				errs.Items = append(errs.Items, e)
				continue
			}
			doc.Blocks = append(doc.Blocks, block)

		case "header":
			var block HeaderBlock
			block.ID = rb.ID
			if err := json.Unmarshal(rb.Data, &block.Data); err != nil {
				errs.Add(path+".data", "invalid header payload: "+err.Error())
				continue
			}
			if e := validateHeader(block.Data, path+".data"); e != nil {
				errs.Items = append(errs.Items, e)
				continue
			}
			doc.Blocks = append(doc.Blocks, block)

		case "list":
			var block ListBlock
			block.ID = rb.ID
			if err := json.Unmarshal(rb.Data, &block.Data); err != nil {
				errs.Add(path+".data", "invalid list payload: "+err.Error())
				continue
			}
			if e := validateList(&block.Data, path+".data"); e != nil {
				errs.Items = append(errs.Items, e)
				continue
			}
			doc.Blocks = append(doc.Blocks, block)

		case "image":
			var block ImageBlock
			block.ID = rb.ID
			if err := json.Unmarshal(rb.Data, &block.Data); err != nil {
				errs.Add(path+".data", "invalid image payload: "+err.Error())
				continue
			}
			if e := validateImage(block.Data, path+".data"); e != nil {
				errs.Items = append(errs.Items, e)
				continue
			}
			if block.Data.File.Type == ImageSrcFile && block.Data.File.FileID != nil {
				addFile(block.Data.File.FileID)
			}
			doc.Blocks = append(doc.Blocks, block)

		case "gallery":
			var block GalleryBlock
			block.ID = rb.ID
			if err := json.Unmarshal(rb.Data, &block.Data); err != nil {
				errs.Add(path+".data", "invalid gallery payload: "+err.Error())
				continue
			}
			if e := validateGallery(&block.Data, path+".data"); e != nil {
				errs.Items = append(errs.Items, e)
				continue
			}
			for _, f := range block.Data.Files {
				if f.Type == ImageSrcFile && f.FileID != nil {
					addFile(f.FileID)
				}
			}
			doc.Blocks = append(doc.Blocks, block)

		default:
			errs.Add(path+".type", "unsupported block type: "+rb.Type)
			continue
		}
	}

	if err := errs.OrNil(); err != nil {
		return ParsedDocument{}, nil, err
	}
	return doc, files, nil
}
