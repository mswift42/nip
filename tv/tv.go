package tv

import (
	"github.com/mswift42/goquery"
	"sync"
	"fmt"
	"strings"
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

type programPage struct {
	doc *iplayerDocument
}

func (pp *programPage) programmes() []*Programme {
	fmt.Println("programpage: ", pp.doc)
	var results []*Programme
	title := pp.doc.doc.Find(".hero-header__title").Text()
	pp.doc.doc.Find(".content-item").Each(func(i int, s *goquery.Selection) {
		results = append(results, newProgrammeFromProgramPage(title, s))
	})
	return results
}

func newProgrammeFromProgramPage(title string, s *goquery.Selection) *Programme {
	subtitle := s.Find(".content-item__title").Text()
	synopsis := s.Find(".content-item__info__secondary > .content-item__description").Text()
	url := s.Find("a").AttrOr("href", "")
	sel := iplayerSelection{s}
	thumbnail := sel.extractThumbnail()
	return &Programme{title, subtitle, synopsis, "", thumbnail, url, 0}
}

func (is *iplayerSelection) extractThumbnail() string {
	set := is.sel.Find(".rs-image > picture > source").AttrOr("srcset", "")
	return strings.Split(set, " ")[0]
}

type Category struct {
	name       string
	programmes []*Programme
}

// TODO - add map for already visited programme sites.
// TODO - implement seenprogramme method
type mainCategoryDocument struct {
	maindoc          *iplayerDocument
	nextdocs         []*iplayerDocument
	selectionresults []*iplayerSelectionResult
}

func (mcd *mainCategoryDocument) programmes() []*Programme {
	var results []*Programme
	for _, i := range mcd.selectionresults {
		if i.prog != nil {
			results = append(results, i.prog)
		}
	}
	return results
}

var seen = make(map[Pager]bool)
var mutex = &sync.Mutex{}
// TODO - replace with map from sync package.
func seenLink(p Pager) bool {
	mutex.Lock()
	if !seen[p] {
		seen[p] = true
		mutex.Unlock()
		return false
	}
	mutex.Unlock()
	return true
}

func (id *iplayerDocument) mainDoc() *iplayerDocument {
	return id
}

func (id *iplayerDocument) nextPages() []Pager {
	var urls []Pager
	id.doc.Find(".page > a").Each(func(i int, s *goquery.Selection) {
		urls = append(urls, BeebURL(s.AttrOr("href", "")))
	})
	return urls
}

func (id *iplayerDocument) programPages() ([]Pager, []*iplayerSelectionResult) {
	var urls []Pager
	urls = append(urls, id.nextPages()...)
	np := collectPages(urls)
	docs := documentsFromResults(np)
	docs = append(docs, id)
	var selres []*iplayerSelectionResult
	for _, i := range docs {
		isel := iplayerSelection{i.doc.Find(".list-item-inner")}
		selres = append(selres, isel.selectionResults()...)
		for _, i := range selres {
			if i.programPage != "" {
				urls = append(urls, BeebURL(i.programPage))
			}
		}
	}
	return urls, selres
}

func documentsFromResults(docres []*iplayerDocumentResult) []*iplayerDocument {
	var results []*iplayerDocument
	for _, i := range docres {
		if i.Error == nil {
			results = append(results, &i.idoc)
		}
	}
	return results
}

func newMainCategory(np NextPager) *mainCategoryDocument {
	var nextdocs []*iplayerDocument
	pp, selres := np.programPages()
	progPages := collectPages(pp)
	for _, i := range progPages {
		if &i.idoc != nil {
			nextdocs = append(nextdocs, &i.idoc)
		}
	}
	return &mainCategoryDocument{np.mainDoc(), nextdocs, selres}
}

func collectPages(urls []Pager) []*iplayerDocumentResult {
	var results []*iplayerDocumentResult
	c := make(chan *iplayerDocumentResult)
	jobs := 0
	for _, i := range urls {
		if !seenLink(i) {
			jobs++
			go func(u Pager) {
				u.loadDocument(c)
			}(i)
		}
	}
	for i := 0; i < jobs; i++ {
		//func (mcd *mainCategoryDocument) programmes() []*Programme {
		results = append(results, <-c)
	}
	return results
}
