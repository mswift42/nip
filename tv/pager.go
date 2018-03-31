package tv

type Pager interface {
	loadDocument(chan<- *iplayerDocumentResult)
}

type NextPager interface {
	collectNextPages([]string) []*iplayerDocumentResult
	collectProgramPages([]string) []*iplayerDocumentResult
}