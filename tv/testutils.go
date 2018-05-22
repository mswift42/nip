package tv

import (
	"bytes"
	"io/ioutil"

	"github.com/mswift42/goquery"
)

type TestHTMLURL string
type TestIplayerDocument struct {
	idoc *iplayerDocument
}

func (thu TestHTMLURL) loadDocument(c chan<- *IplayerDocumentResult) {
	file, err := ioutil.ReadFile(string(thu))
	if err != nil {
		c <- &IplayerDocumentResult{iplayerDocument{}, err}
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(file))
	if err != nil {
		c <- &IplayerDocumentResult{iplayerDocument{}, err}
	}
	idoc := iplayerDocument{doc, thu}
	c <- &IplayerDocumentResult{idoc, nil}
}

func (tid *TestIplayerDocument) nextPages() []Pager {
	var urls []Pager
	tid.idoc.doc.Find(".page > a").Each(func(i int, s *goquery.Selection) {
		urls = append(urls, TestHTMLURL(s.AttrOr("href", "")))
	})
	return urls
}

func (tid *TestIplayerDocument) programPages(selres []*iplayerSelectionResult) []Pager {
	var urls []Pager
	for _, i := range selres {
		if i.programPage != "" {
			urls = append(urls, TestHTMLURL(i.programPage))
		}
	}
	return urls
}

func (tid *TestIplayerDocument) mainDoc() *iplayerDocument {
	return tid.idoc
}

func documentLoader(url string) *iplayerDocument {
	thu := TestHTMLURL(url)
	c := make(chan *IplayerDocumentResult)
	go thu.loadDocument(c)
	idr := <-c
	if idr.Error != nil {
		panic(idr.Error)
	}
	return &idr.Idoc
}

func RemoteDocumentLoader(url string) *iplayerDocument {
	bu := BeebURL(url)
	c := make(chan *IplayerDocumentResult)
	go bu.loadDocument(c)
	idr := <-c
	if idr.Error != nil {
		panic(idr.Error)
	}
	return &idr.Idoc
}

func contains(progs []*Programme, url string) bool {
	for _, i := range progs {
		if i.URL == url {
			return true
		}
	}
	return false
}

func findProgramme(progs []*Programme, url string) *Programme {
	for _, i := range progs {
		if i.URL == url {
			return i
		}
	}
	return nil
}