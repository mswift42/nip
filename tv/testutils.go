package tv

import (
	"bytes"
	"io/ioutil"

	"github.com/mswift42/goquery"
)

type TestHTMLURL string
type TestIplayerDocument struct {
	idoc iplayerDocument
}

func (thu TestHTMLURL) loadDocument(c chan<- *iplayerDocumentResult) {
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

func (tid TestIplayerDocument) nextPages() []interface{} {
	var urls []interface{}
	tid.idoc.doc.Find(".page > a").Each(func(i int, s *goquery.Selection) {
		urls = append(urls, TestHTMLURL(s.AttrOr("href", "")))
	})
	return urls
}

func (thu TestHTMLURL) collectPages(urls []string) []*iplayerDocumentResult {
	var results []*iplayerDocumentResult
	c := make(chan *iplayerDocumentResult)
	for _, i := range urls {
		go func(s string) {
			thu := TestHTMLURL(s)
			thu.loadDocument(c)
		}(i)
	}
	for i := 0; i < len(urls); i++ {
		results = append(results, <-c)
	}
	return results
}
