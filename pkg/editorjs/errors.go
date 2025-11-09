package editorjs

import (
	"fmt"
	"net/url"
	"strings"
)

type ValidationError struct {
	Path   string
	Reason string
}

func (e ValidationError) Error() string { return fmt.Sprintf("%s: %s", e.Path, e.Reason) }

type MultiError struct {
	Items []error
}

func (m *MultiError) Error() string {
	if len(m.Items) == 0 {
		return ""
	}
	var b strings.Builder
	for i, e := range m.Items {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(e.Error())
	}
	return b.String()
}

func (m *MultiError) Add(path, reason string) {
	m.Items = append(m.Items, ValidationError{Path: path, Reason: reason})
}

func (m *MultiError) OrNil() error {
	if len(m.Items) == 0 {
		return nil
	}
	return m
}

func isHTTPURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return u.Host != ""
}

func isValidURL(s string) bool {
	_, err := url.Parse(s)
	return err == nil
}
