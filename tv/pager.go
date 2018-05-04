package tv

type Pager interface {
	loadDocument(chan<- *IplayerDocumentResult)
}

type NextPager interface {
	mainDoc() *iplayerDocument
	nextPages() []Pager
	programPages([]*iplayerDocument) ([]Pager, []*iplayerSelectionResult)
}