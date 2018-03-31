package tv

type Pager interface {
	loadDocument(chan<- *iplayerDocumentResult)
}

type NextPager interface {
	collectPages([]string) []*iplayerDocumentResult
}