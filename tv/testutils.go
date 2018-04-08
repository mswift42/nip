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

func (tid TestIplayerDocument) programPages() []interface{} {
	var urls []interface{}
	isel := iplayerSelection{tid.idoc.doc.Find(".list-item-inner")}
	selres := isel.selectionResults()
	for _, i := range selres {
		if i.programPage != "" {
			urls = append(urls, TestHTMLURL(i.programPage))
		}
	}
	return urls
}

