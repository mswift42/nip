package tv

type Pager interface {
	loadDocument(chan<- *iplayerDocumentResult)
}

type NextPager interface {
	mainDoc() *iplayerDocument
	nextPages() []interface{}
	programPages() []interface{}
}