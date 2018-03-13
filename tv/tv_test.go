package tv

import "testing"

func TestLoadingDocument(t *testing.T) {
	url := TestHtmlUrl("testhtml/food1.html")
	idr := url.loadDocument()
	if idr.Error != nil {
		t.Error("Expected error to be nil", idr.Error)
	}
	if idr.idoc.doc == nil {
		t.Error("Expected idoc not to be nil", idr.idoc)
	}
	url = TestHtmlUrl("testhtml/films1.html")
	idr = url.loadDocument()
	if idr.Error != nil {
		t.Error("Expected error to be nil: ", idr.Error)
	}
	if idr.idoc.doc == nil {
		t.Error("Expected idoc not to be nil: ", idr.idoc)
	}
}
func TestIplayerSelectionResults(t *testing.T) {
	url := TestHtmlUrl("testhtml/films1.html")
	idr := url.loadDocument()
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

func TestMainCategoryDocumentNextPages(t *testing.T) {
	url := TestHtmlUrl("testhtml/films1.html")
	idr := url.loadDocument()
	if idr.Error != nil {
		t.Error("Expected no error loading document, got: ", idr.Error)
	}
	var emptydoc []*iplayerDocument
	mcd := mainCategoryDocument{&idr.idoc, emptydoc}
	np := mcd.nextPages()
	if len(np) != 1 {
		t.Error("Expected length of nextpages to be 1, got: ", len(np))
	}
	if np[0] != "films2.html" {
		t.Error("Expected url of first nextPage to be films2.html, got: ", np[0])
	}
}
