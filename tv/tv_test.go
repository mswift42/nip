package tv

import "testing"

func TestLoadingDocument(t *testing.T) {
	url := TestHtmlUrl("testhtml/food1.html")
	c := make(chan *iplayerDocumentResult)
	url.loadDocument(c)
	idr := <-c
	if idr.Error != nil {
		t.Error("Expected error to be nil", idr.Error)
	}
	if idr.idoc.doc == nil {
		t.Error("Expected idoc not to be nil", idr.idoc)
	}
	url = TestHtmlUrl("testhtml/films1.html")
	url.loadDocument(c)
	idr = <-c
	if idr.Error != nil {
		t.Error("Expected error to be nil: ", idr.Error)
	}
	if idr.idoc.doc == nil {
		t.Error("Expected idoc not to be nil: ", idr.idoc)
	}
}
func TestIplayerSelectionResults(t *testing.T) {
	url := TestHtmlUrl("testhtml/films1.html")
	c := make(chan *iplayerDocumentResult)
	url.loadDocument(c)
	idr := <-c
	sel := iplayerSelection{idr.idoc.doc.Find(".list-item-inner")}
	selres := sel.selectionResults()
	if len(selres) != 20 {
		t.Error("Expected length of selectionresults to equal: ", len(selres))
	}
	progpage := selres[0]
	if progpage.prog != nil {
		t.Error("Expected proramme to be nil: ", progpage.prog)
	}
	if progpage.programPage != "adam_curtis.html" {
		t.Error("Expected program Page to be 'adam_curtis.html' not: ", progpage.programPage)
	}
	if selres[1].prog.Title != "A Hijacking" {
		t.Error("Expected second programme title to be 'A Hijacking', got: ", selres[1].prog.Title)
	}
	if selres[1].programPage != "" {
		t.Error("Expected second programPage to be an empty string, got: ", selres[1].programPage)
	}
}

//func TestNewTestMainCategory(t *testing.T) {
//	url := TestHtmlUrl("testhtml/films1.html")
//	nmc := url.newMainCategory()
//	if nmc.maindoc == nil {
//		t.Error("Expected maindocument to not be nil, got: ", nmc.maindoc)
//	}
//	if len(nmc.nextdocs) != 1 {
//		t.Error("Expected length of nextdocs to be 1, got: ", len(nmc.nextdocs))
//	}
//	sel := iplayerSelection{nmc.nextdocs[0].doc.Find(".list-item-inner")}
//	selres := sel.selectionResults()
//	if len(selres) != 4 {
//		t.Error("Expected length of selectionresutls of films2.html to be 4, got: ", len(selres))
//	}
//}

