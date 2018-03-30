package tv

import (
	"bytes"
	"io/ioutil"

	"github.com/mswift42/goquery"
)

type TestHtmlUrl string

func (thu TestHtmlUrl) newMainCategory() *mainCategoryDocument {
	var results []*iplayerDocument
	maindocres := thu.loadDocument()
	if maindocres.Error != nil {
		return &mainCategoryDocument{ nil, results}
	}
	np := maindocres.idoc.nextPages()
	npresults := collectNextPages(np)
	for _, i := range npresults {
		if i.Error == nil {
			results = append(results, &i.idoc)
		}
	}
	return &mainCategoryDocument{&maindocres.idoc, results}
}


func (thu TestHtmlUrl) loadDocument(c chan<- *iplayerDocumentResult) *iplayerDocumentResult {
	file, err := ioutil.ReadFile(string(thu))
	if err != nil {
		c <-  &iplayerDocumentResult{iplayerDocument{}, err}
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(file))
	if err != nil {
		c <- return &iplayerDocumentResult{iplayerDocument{}, err}
	}
	idoc := iplayerDocument{doc}
	c <-  &iplayerDocumentResult{idoc, nil}
}
