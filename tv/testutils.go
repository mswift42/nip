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

func (thu TestHtmlUrl) newMainCategory() []*iplayerDocumentResult {
	var results []*iplayerDocumentResult
	maindocres := thu.loadDocument()
	if maindocres.Error != nil {
		return results
	}
	np := maindocres.idoc.nextPages()
	results = collectNextPages(np)
	return results
}

func  collectNextPages(urls []string) []*iplayerDocumentResult {
	var results []*iplayerDocumentResult
	c := make(chan *iplayerDocumentResult)
	for _, i := range urls {
		go func(u string) {
			th := TestHtmlUrl(u)
			idr := th.loadDocument()
			c <- idr
		}(i)
	}
	close (c)
	for range c {
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
