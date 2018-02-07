package tv

type Searcher interface {
	loadDocument() *iplayerDocumentResult
}
