package tv

import (
	"github.com/mswift42/goquery"
)

type BeebURL string

func (bu BeebURL) loadDocument(c chan<- *iplayerDocumentResult) {
	doc, err := goquery.NewDocument(string(bu))
	if err != nil {
		c <- &iplayerDocumentResult{iplayerDocument{}, err}
	}
	idoc := iplayerDocument{doc}
	c <- &iplayerDocumentResult{idoc, nil}
}

type iplayerSelection struct {
	sel *goquery.Selection
}

func newIplayerSelection(sel *goquery.Selection) *iplayerSelection {
	return &iplayerSelection{sel}
}

// iplayerSelectionResult has either an iplayer programme,
// or if it has a "more Programmes available" notice, the link to its Programme Page.
type iplayerSelectionResult struct {
	prog        *Programme
	programPage string
}

func (is *iplayerSelection) selectionResults() []*iplayerSelectionResult {
	var res []*iplayerSelectionResult
	is.sel.Each(func(i int, selection *goquery.Selection) {
		isel := newIplayerSelection(selection)
		page := isel.programmeSite()
		if page != "" {
			res = append(res, &iplayerSelectionResult{nil, page})
		} else {
			res = append(res,
				&iplayerSelectionResult{isel.programme(), ""})
		}
	})
	return res
}

func (is *iplayerSelection) programmeSite() string {
	return is.sel.Find(".view-more-container").AttrOr("href", "")
}

func (is *iplayerSelection) programme() *Programme {
	title := is.title()
	subtitle := is.subtitle()
	synopsis := is.synopsis()
	url := is.url()
	thumbnail := is.thumbNail()
	pid := is.pid()
	return &Programme{
		Title:     title,
		Subtitle:  subtitle,
		Synopsis:  synopsis,
		PID:       pid,
		Thumbnail: thumbnail,
		URL:       url,
		Index:     0,
	}
}

func (is *iplayerSelection) title() string {
	return is.sel.Find(".secondary > .title").Text()
}

func (is *iplayerSelection) subtitle() string {
	return is.sel.Find(".secondary > .subtitle").Text()
}

func (is *iplayerSelection) synopsis() string {
	return is.sel.Find(".synopsis").Text()
}

func (is *iplayerSelection) url() string {
	return is.sel.Find("a").AttrOr("href", "")
}

func (is *iplayerSelection) thumbNail() string {
	return is.sel.Find(".rs-image > picture > source").AttrOr("srcset", "")
}

func (is *iplayerSelection) pid() string {
	pid := is.sel.AttrOr("data-ip-id", "")
	if pid != "" {
		return pid
	}
	return is.sel.Find(".list-item-inner > a").AttrOr("data-episode-id", "")
}

// Programme represents an Iplayer TV programme. It consists of
// the programme's title, subtitle, a short programme description,
// The Iplayer Programme ID, the url to its thumbnail, the url
// to the programme's website and a unique index.
type Programme struct {
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Synopsis  string `json:"synopsis"`
	PID       string `json:"pid"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"url"`
	Index     int    `json:"index"`
}

type iplayerDocument struct {
	doc *goquery.Document
}

type iplayerDocumentResult struct {
	idoc  iplayerDocument
	Error error
}
type Category struct {
	name       string
	programmes []*Programme
}

type mainCategoryDocument struct {
	maindoc  *iplayerDocument
	nextdocs []*iplayerDocument
}

func (mcd *mainCategoryDocument) nextPages() []string {
	var url []string
	mcd.maindoc.doc.Find(".page > a").Each(func(i int, s *goquery.Selection) {
		url = append(url, s.AttrOr("href", ""))
	})
	return url
}

func (id *iplayerDocument) mainDoc() *iplayerDocument {
	return id
}

func (id *iplayerDocument) nextPages() []Pager  {
	var urls []Pager
	id.doc.Find(".page > a").Each(func(i int, s *goquery.Selection) {
		urls = append(urls, BeebURL(s.AttrOr("href", "")))
	})
	return urls
}

func (id *iplayerDocument) programPages() []Pager {
	var urls []Pager
	isel := iplayerSelection{id.doc.Find(".list-item-inner")}
	selres := isel.selectionResults()
	for _, i := range selres {
		if i.programPage != "" {
			urls = append(urls, BeebURL(i.programPage))
		}
	}
	return urls
}

func (bu BeebURL) collectPages(urls []string) []*iplayerDocumentResult {
	var results []*iplayerDocumentResult
	c := make(chan *iplayerDocumentResult)
	for _, i := range urls {
		go func(s string) {
			bu := BeebURL(s)
			bu.loadDocument(c)
		}(i)
	}
	for i := 0; i < len(urls); i++ {
		results = append(results, <-c)
	}
	return results
}

func newMainCategory(np NextPager) *mainCategoryDocument {
	var pages []*iplayerDocument
	nextPages := collectPages(np.nextPages())
	progPages := collectPages(np.programPages())
	for _, i := range nextPages {
		if i.Error == nil {
			pages = append(pages, &i.idoc)
		}
	}
	for _, i := range progPages {
		if i.Error == nil {
			pages = append(pages, &i.idoc)
		}
	}
	return &mainCategoryDocument{np.mainDoc(), pages}
}

func collectPages(urls []Pager) []*iplayerDocumentResult {
	var results []*iplayerDocumentResult
	c := make(chan *iplayerDocumentResult)
	for _, i := range urls {
		go func(u Pager) {
				u.loadDocument(c)
		}(i)
	}
	for i := 0; i < len(urls); i++ {
		results = append(results, <-c)
	}
	return results
}

