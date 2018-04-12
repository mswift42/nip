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

func (tid *TestIplayerDocument) nextPages() []Pager {
	var urls []Pager
	tid.idoc.doc.Find(".page > a").Each(func(i int, s *goquery.Selection) {
		urls = append(urls, TestHTMLURL(s.AttrOr("href", "")))
	})
	return urls
}

func (tid *TestIplayerDocument) programPages() ([]Pager, []*iplayerSelectionResult) {
	var urls []Pager
	urls = append(urls, tid.nextPages()...)
	np := collectPages(urls)
	docs := documentsFromResults(np)
	docs = append(docs, &tid.idoc)
	var selres []*iplayerSelectionResult
	for _, i := range docs {
		isel := iplayerSelection{i.doc.Find(".list-item-inner")}
		selres = append(selres, isel.selectionResults()...)
		for _, i := range selres {
			if i.programPage != "" {
				urls = append(urls, TestHTMLURL(i.programPage))
			}
		}
	}
	return urls, selres
}

func (tid *TestIplayerDocument) mainDoc() *iplayerDocument {
	return &tid.idoc
}
