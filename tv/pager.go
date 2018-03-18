package tv

type Pager interface {
	loadDocument() *iplayerDocumentResult
}

type NextPager interface {
	collectNextPages([]string) []*iplayerDocumentResult
}