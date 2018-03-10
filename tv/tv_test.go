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
	sel := iplayerSelection {idr.idoc.doc.Find(".list-item-inner")}
	selres := sel.selectionResults()
	if len(selres) != 20 {
		t.Error("Expected length of selectionresults to equal: ", len(selres))
	}
	progpage := selres[0]
	if progpage.prog != nil {
		t.Error("Expected proramme to be nil: ", progpage.prog)
	}
	if progpage.programPage != BeebUrl("adam_curtis.html") {
		t.Error("Expected program Page to be 'adam_curtis.html' not: ", progpage.programPage)
	}
}
