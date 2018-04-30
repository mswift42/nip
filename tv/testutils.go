package tv

import (
	"bytes"
	"io/ioutil"

	"github.com/mswift42/goquery"
	"fmt"
)

type TestHTMLURL string
type TestIplayerDocument struct {
	idoc *iplayerDocument
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
	fmt.Println("Printing next Pages: ", tid.nextPages())
	urls = append(urls, tid.nextPages()...)
	np := collectPages(urls)
	docs := []*iplayerDocument{tid.idoc}
	docs = append(docs, documentsFromResults(np)...)
	fmt.Println(docs)
	var selres []*iplayerSelectionResult
	for _, i := range docs {
		isel := iplayerSelection{i.doc.Find(".list-item-inner")}
		selres = append(selres, isel.selectionResults()...)
		fmt.Println(selres)
	}
	for _, i := range selres {
		if i.programPage != "" {
			fmt.Println("Printing Program Page: ", i.programPage)
			urls = append(urls, TestHTMLURL(i.programPage))
		}
	}
	fmt.Println(urls)
	return urls, selres
}

func (tid *TestIplayerDocument) mainDoc() *iplayerDocument {
	return tid.idoc
}

func documentLoader(url string) *iplayerDocument {
	thu := TestHTMLURL(url)
	c := make(chan *iplayerDocumentResult)
	go thu.loadDocument(c)
	idr := <-c
	if idr.Error != nil {
		panic(idr.Error)
	}
	return &idr.idoc
}

func contains(progs[]*Programme, url string) bool {
	for _, i := range progs {
		if i.URL == url {
			return true
		}
	}
	return false
}
