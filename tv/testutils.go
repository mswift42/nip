package tv

import (
	"bytes"
	"io/ioutil"

	"github.com/mswift42/goquery"
)

type TestHtmlUrl string

type testMainCategoryDocument struct {
	ip *iplayerDocument
	nextPages []*iplayerDocument
}

func (tmcd testMainCategoryDocument) collectNextPages(urls []string) []*iplayerDocumentResult {
	var results []*iplayerDocumentResult
	c := make(chan *iplayerDocumentResult)
	for _, i := range urls {
		go func(u url) {
			th := TestHtmlUrl(u)
			idr := th.loadDocument()
			c <- idr
		}(i)
	}
	for i := 0; i < len(urls); i++ {
		results = append(results, <-c)
	}
	return results
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
