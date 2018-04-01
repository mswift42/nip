package tv

import (
	"bytes"
	"io/ioutil"

	"github.com/mswift42/goquery"
)

type TestHtmlUrl string

func (thu TestHtmlUrl) loadDocument(c chan<- *iplayerDocumentResult) {
	file, err := ioutil.ReadFile(string(thu))
	if err != nil {
		c <- &iplayerDocumentResult{iplayerDocument{}, err}
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(file))
	if err != nil {
		c <- &iplayerDocumentResult{iplayerDocument{}, err}
	}
	idoc := iplayerDocument{doc}
	c <- &iplayerDocumentResult{idoc, nil}
}

func (thu TestHtmlUrl) collectPages(urls []string) []*iplayerDocumentResult {
	var results []*iplayerDocumentResult
	c := make(chan *iplayerDocumentResult)
	for _, i := range urls {
		go func(s string) {
			thu := TestHtmlUrl(s)
			thu.loadDocument(c)
		}(i)
	}
	for i := 0; i < len(urls); i++ {
		results = append(results, <-c)
	}
	return results
}
