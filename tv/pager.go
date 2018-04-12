package tv

type Pager interface {
	loadDocument(chan<- *iplayerDocumentResult)
}

type NextPager interface {
	mainDoc() *iplayerDocument
	nextPages() []Pager
	programPages() ([]Pager, []*iplayerSelectionResult)
}