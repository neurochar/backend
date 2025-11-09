package editorjs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UnixMillis time.Time

func (u *UnixMillis) UnmarshalJSON(b []byte) error {
	var ms int64
	if err := json.Unmarshal(b, &ms); err != nil {
		return fmt.Errorf("time must be unix millis: %w", err)
	}
	t := time.UnixMilli(ms)
	*u = UnixMillis(t)
	return nil
}

func (u UnixMillis) Time() time.Time { return time.Time(u) }

type rawDocument struct {
	Time    UnixMillis `json:"time"`
	Blocks  []rawBlock `json:"blocks"`
	Version string     `json:"version"`
}

type rawBlock struct {
	ID   string          `json:"id"`
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ParsedDocument struct {
	Time    time.Time
	Version string
	Blocks  []AnyBlock
}

type AnyBlock interface {
	Kind() string
	BlockID() string
}

type ParagraphAlignment string

const (
	AlignLeft    ParagraphAlignment = "left"
	AlignCenter  ParagraphAlignment = "center"
	AlignRight   ParagraphAlignment = "right"
	AlignJustify ParagraphAlignment = "justify"
)

type ParagraphData struct {
	Text      string             `json:"text"`
	Alignment ParagraphAlignment `json:"alignment"`
}

type ParagraphBlock struct {
	ID   string        `json:"id"`
	Data ParagraphData `json:"data"`
}

func (b ParagraphBlock) Kind() string    { return "paragraph" }
func (b ParagraphBlock) BlockID() string { return b.ID }

type HeaderData struct {
	Text  string `json:"text"`
	Level int    `json:"level"`
}

type HeaderBlock struct {
	ID   string     `json:"id"`
	Data HeaderData `json:"data"`
}

func (b HeaderBlock) Kind() string    { return "header" }
func (b HeaderBlock) BlockID() string { return b.ID }

type ListStyle string

const (
	ListUnordered ListStyle = "unordered"
	ListOrdered   ListStyle = "ordered"
)

type ListItem struct {
	Content string         `json:"content"`
	Meta    map[string]any `json:"meta"`
	Items   []ListItem     `json:"items"`
}

type ListData struct {
	Style ListStyle      `json:"style"`
	Meta  map[string]any `json:"meta"`
	Items []ListItem     `json:"items"`
}

type ListBlock struct {
	ID   string   `json:"id"`
	Data ListData `json:"data"`
}

func (b ListBlock) Kind() string    { return "list" }
func (b ListBlock) BlockID() string { return b.ID }

type ImageSourceType string

const (
	ImageSrcURL  ImageSourceType = "url"
	ImageSrcFile ImageSourceType = "file"
)

type ImageFile struct {
	URL        string          `json:"url"`
	Type       ImageSourceType `json:"type"`
	FileID     *uuid.UUID      `json:"fileID,omitempty"`
	Filename   string          `json:"filename,omitempty"`
	FileTarget string          `json:"fileTarget,omitempty"`
}

type ImageData struct {
	File    ImageFile `json:"file"`
	Caption string    `json:"caption"`
}

type ImageBlock struct {
	ID   string    `json:"id"`
	Data ImageData `json:"data"`
}

func (b ImageBlock) Kind() string    { return "image" }
func (b ImageBlock) BlockID() string { return b.ID }

type GalleryStyle string

const (
	GalleryGrid   GalleryStyle = "grid"
	GallerySlider GalleryStyle = "slider"
)

type GalleryFile struct {
	URL        string          `json:"url"`
	Type       ImageSourceType `json:"type"`
	FileID     *uuid.UUID      `json:"fileID,omitempty"`
	Filename   string          `json:"filename,omitempty"`
	FileTarget string          `json:"fileTarget,omitempty"`
	Caption    string          `json:"caption"`
}

type GalleryData struct {
	Files   []GalleryFile `json:"files"`
	Caption string        `json:"caption"`
	Style   GalleryStyle  `json:"style"`
}

type GalleryBlock struct {
	ID   string      `json:"id"`
	Data GalleryData `json:"data"`
}

func (b GalleryBlock) Kind() string    { return "gallery" }
func (b GalleryBlock) BlockID() string { return b.ID }
