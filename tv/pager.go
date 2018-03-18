package tv

type Pager interface {
	loadDocument() *iplayerDocumentResult
	collectNextPages([]string) []*iplayerDocumentResult
}
