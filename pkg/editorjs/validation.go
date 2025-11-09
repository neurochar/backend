package editorjs

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/net/html"
)

func validateParagraph(d *ParagraphData, path string) error {
	var m MultiError
	safe, err := sanitizeHTML(d.Text)
	if err != nil {
		m.Add(path+".text", "invalid HTML: "+err.Error())
	} else {
		d.Text = safe
	}

	switch d.Alignment {
	case AlignLeft, AlignCenter, AlignRight, AlignJustify:
	default:
		m.Add(path+".alignment", "must be one of: left, center, right, justify")
	}
	return m.OrNil()
}

func validateHeader(d HeaderData, path string) error {
	var m MultiError
	if strings.TrimSpace(d.Text) == "" {
		m.Add(path+".text", "must be non-empty")
	}
	if d.Level < 1 || d.Level > 6 {
		m.Add(path+".level", "must be between 1 and 6")
	}
	return m.OrNil()
}

func validateList(d *ListData, path string) error {
	var m MultiError
	switch d.Style {
	case ListUnordered, ListOrdered:
	default:
		m.Add(path+".style", "must be one of: unordered, ordered")
	}
	if len(d.Items) == 0 {
		m.Add(path+".items", "must contain at least 1 item")
	} else {
		sanitizeAndValidateListItems(d.Items, path+".items", 1, &m)
	}
	return m.OrNil()
}

func sanitizeAndValidateListItems(items []ListItem, path string, depth int, m *MultiError) {
	if depth > 10 {
		m.Add(path, "list nesting depth must not exceed 10")
		return
	}

	for i := range items {
		p := fmt.Sprintf("%s[%d]", path, i)
		safe, err := sanitizeHTML(items[i].Content)
		if err != nil {
			m.Add(p+".content", "invalid HTML: "+err.Error())
		} else {
			items[i].Content = safe
		}
		if strings.TrimSpace(stripAllHTML(items[i].Content)) == "" {
			m.Add(p+".content", "must be non-empty")
		}
		if len(items[i].Items) > 0 {
			sanitizeAndValidateListItems(items[i].Items, p+".items", depth+1, m)
		}
	}
}

func validateImage(d ImageData, path string) error {
	var m MultiError
	validateImageFile(d.File, path+".file", &m)
	return m.OrNil()
}

func validateGallery(d *GalleryData, path string) error {
	var m MultiError
	switch d.Style {
	case GalleryGrid, GallerySlider:
	default:
		m.Add(path+".style", "must be one of: grid, slider")
	}
	if len(d.Files) == 0 {
		m.Add(path+".files", "must contain at least 1 file")
	} else {
		for i := range d.Files {
			validateImageFile(
				ImageFile{
					URL:        d.Files[i].URL,
					Type:       d.Files[i].Type,
					FileID:     d.Files[i].FileID,
					Filename:   d.Files[i].Filename,
					FileTarget: d.Files[i].FileTarget,
				},
				fmt.Sprintf("%s.files[%d]", path, i),
				&m,
			)
		}
	}
	return m.OrNil()
}

func validateImageFile(f ImageFile, path string, m *MultiError) {
	if !isHTTPURL(f.URL) {
		m.Add(path+".url", "must be valid http/https URL")
	}
	switch f.Type {
	case ImageSrcURL:
	case ImageSrcFile:
		if f.FileID == nil || *f.FileID == uuid.Nil {
			m.Add(path+".fileID", "must be valid UUID for type=file")
		}
		if strings.TrimSpace(f.Filename) == "" {
			m.Add(path+".filename", "required for type=file")
		}
		if strings.TrimSpace(f.FileTarget) == "" {
			m.Add(path+".fileTarget", "required for type=file")
		}
	default:
		m.Add(path+".type", "must be one of: url, file")
	}
}

var (
	allowedTags       = map[string]struct{}{"a": {}, "b": {}, "i": {}, "u": {}}
	allowedInlineTags = map[string]struct{}{"br": {}}
	allowedTargets    = map[string]struct{}{"_blank": {}, "_self": {}, "_parent": {}, "_top": {}}
	relTokenRe        = regexp.MustCompile(`^[a-zA-Z-]+$`)
)

func sanitizeHTML(input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", nil
	}
	nodes, err := html.ParseFragment(strings.NewReader(input), nil)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	for _, n := range nodes {
		if err := renderSanitized(&b, n); err != nil {
			return "", err
		}
	}
	return b.String(), nil
}

func renderSanitized(b *strings.Builder, n *html.Node) error {
	switch n.Type {
	case html.TextNode:
		b.WriteString(html.EscapeString(n.Data))
	case html.ElementNode:
		name := strings.ToLower(n.Data)

		if _, ok := allowedTags[name]; ok {
			b.WriteByte('<')
			b.WriteString(name)
			if name == "a" {
				for _, a := range filterAnchorAttrs(n.Attr) {
					b.WriteByte(' ')
					b.WriteString(a.Key)
					b.WriteString(`="`)
					b.WriteString(html.EscapeString(a.Val))
					b.WriteByte('"')
				}
			}
			b.WriteByte('>')
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if err := renderSanitized(b, c); err != nil {
					return err
				}
			}
			b.WriteString("</")
			b.WriteString(name)
			b.WriteByte('>')
		} else if _, ok := allowedInlineTags[name]; ok {
			b.WriteByte('<')
			b.WriteString(name)
			b.WriteString("/>")
		} else {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if err := renderSanitized(b, c); err != nil {
					return err
				}
			}
		}
	default:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := renderSanitized(b, c); err != nil {
				return err
			}
		}
	}
	return nil
}

func filterAnchorAttrs(attrs []html.Attribute) []html.Attribute {
	var out []html.Attribute
	var href, target, rel *string
	for _, a := range attrs {
		key := strings.ToLower(a.Key)
		val := strings.TrimSpace(a.Val)
		switch key {
		case "href":
			if isValidURL(val) {
				href = &val
			}
		case "target":
			if _, ok := allowedTargets[strings.ToLower(val)]; ok {
				target = &val
			}
		case "rel":
			tokens := strings.Fields(val)
			var filtered []string
			for _, t := range tokens {
				if relTokenRe.MatchString(t) {
					switch strings.ToLower(t) {
					case "noopener", "noreferrer", "nofollow", "ugc", "external":
						filtered = append(filtered, t)
					}
				}
			}
			if len(filtered) > 0 {
				v := strings.Join(filtered, " ")
				rel = &v
			}
		}
	}
	if href != nil {
		out = append(out, html.Attribute{Key: "href", Val: *href})
	}
	if target != nil {
		out = append(out, html.Attribute{Key: "target", Val: *target})
	}
	if rel != nil {
		out = append(out, html.Attribute{Key: "rel", Val: *rel})
	}
	return out
}

func stripAllHTML(s string) string {
	var out strings.Builder
	inTag := false
	for _, r := range s {
		switch r {
		case '<':
			inTag = true
		case '>':
			inTag = false
		default:
			if !inTag {
				out.WriteRune(r)
			}
		}
	}
	return out.String()
}
