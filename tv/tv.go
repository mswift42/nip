package tv

import "github.com/mswift42/goquery"

type BeebUrl string

func (bu BeebUrl) loadDocument() *iplayerDocumentResult {
	doc, err := goquery.NewDocument(string(bu))
	if err != nil {
		return &iplayerDocumentResult{iplayerDocument{}, err}
	}
	idoc := iplayerDocument{doc}
	return &iplayerDocumentResult{idoc, nil}
}

// Programme represents an Iplayer TV programme. It consists of
// the programme's title, subtitle, a short programme description,
// The Iplayer Programme ID, the url to its thumbnail, the url
// to the programme's website and a unique index.
type Programme struct {
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Synopsis  string `json:"synopsis"`
	Pid       string `json:"pid"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"url"`
	Index     int    `json:"index"`
}

type iplayerDocument struct {
	idoc *goquery.Document
}

type iplayerDocumentResult struct {
	idoc  iplayerDocument
	Error error
}

type Category struct {
	name string
	docs []*iplayerDocument
}
