package tv

// Pager is the interface that wraps the loadDocument method.
//
// loadDocument sends the result of generating a new iplayerDocument
// from either a remote Url, or the path to a local html file to
// a provided channel of type *IplayerDocumentResult.
type Pager interface {
	loadDocument(chan<- *IplayerDocumentResult)
}

// NextPager is the interface that wraps the mainDoc, nextPages and programPages
// methods.
//
// The mainDoc method returns the root document of a Category, nextPages returns
// a slice of category pages 2 - n, and programPages returns the urls for all
// programmes on pages 1 - n, that have more than one episode available.
type NextPager interface {
	mainDoc() *iplayerDocument
	nextPages() []Pager
	programPages([]*iplayerSelectionResult) []Pager
}
