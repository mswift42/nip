package tv

import (
	"bytes"
	"io/ioutil"

	"github.com/mswift42/goquery"
)

type TestHtmlUrl string

type testMainCategoryDocument struct {
	ip *iplayerDocument
	NextPages []string
}

func (thu TestHtmlUrl) loadDocument() *iplayerDocumentResult {
	file, err := ioutil.ReadFile(string(thu))
	if err != nil {
		return &iplayerDocumentResult{iplayerDocument{}, err}
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(file))
	if err != nil {
		return &iplayerDocumentResult{iplayerDocument{}, err}
	}
	idoc := iplayerDocument{doc}
	return &iplayerDocumentResult{idoc, nil}
}
