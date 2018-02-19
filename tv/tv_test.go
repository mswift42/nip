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
}

func TestIplayerSelections(t *testing.T) {
	url := TestHtmlUrl("testhtml/food1.html")
	idr := url.loadDocument()
	is := idr.idoc.iplayerSelections()
	is1 := is.sel.Nodes[0]
	if is1.Data != "li" {
		t.Error("expected is1.Data to be li", is1.Data)
	}
}
