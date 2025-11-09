package emailing

type Message struct {
	To            string
	Subject       string
	TextHtml      string
	TextPlain     string
	AutoTextPlain bool
	Files         []MessageFile
}

type MessageFile struct {
	Name    string
	Data    []byte
	Headers map[string][]string
}
