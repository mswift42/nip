package tv

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/mswift42/goquery"
)

// BeebURL is the url for an iplayer web site.
type BeebURL string

const bbcprefix = "https://bbc.co.uk"

func (bu BeebURL) loadDocument(c chan<- *IplayerDocumentResult) {
	var url string
	if strings.HasPrefix(string(bu), "/iplayer/") {
		url = "https://www.bbc.co.uk" + string(bu)
	} else {
		url = string(bu)
	}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		c <- &IplayerDocumentResult{iplayerDocument{}, err}
	}
	idoc := iplayerDocument{doc, bu}
	c <- &IplayerDocumentResult{idoc, nil}
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
	return is.sel.Find(".lnk").AttrOr("href", "")
}

func (is *iplayerSelection) programme() *Programme {
	title := is.title()
	subtitle := is.subtitle()
	return newProgramme(title, subtitle, is)
}

func (is *iplayerSelection) title() string {
	return is.sel.Find(".content-item__title").Text()
}

func (is *iplayerSelection) subtitle() string {
	return is.sel.Find(".content-item__info__primary > " +
		".content-item__description").Text()
}

func (is *iplayerSelection) synopsis() string {
	return is.sel.Find(".content-item__info__secondary " +
		"> .content-item__description").Text()
}

func (is *iplayerSelection) url() string {
	return is.sel.Find("a").AttrOr("href", "")
}

func (is *iplayerSelection) thumbnail() string {
	set := is.sel.Find(".rs-image > picture > source").AttrOr("srcset", "")
	return strings.Split(set, " ")[0]
}

func (is *iplayerSelection) available() string {
	return is.sel.Find(".content-item__sublabels > span").Last().Text()
}

func (is *iplayerSelection) duration() string {
	return is.sel.Find(".content-item__sublabels > span").First().Text()
}

// Programme represents an Iplayer TV programme. It consists of
// the programme's title, subtitle, a short programme description,
// The Iplayer Programme ID, the url to its thumbnail, the url
// to the programme's website and a unique index.
type Programme struct {
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Synopsis  string `json:"synopsis"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"url"`
	Index     int    `json:"index"`
	Available string `json:"available"`
	Duration  string `json:"duration"`
}

func (p *Programme) String() string {
	return fmt.Sprintf("%d: %s %s, %s, %s\n",
		p.Index, p.Title, p.Subtitle, p.Available, p.Duration)
}

type iplayerDocument struct {
	doc *goquery.Document
	url Pager
}

func (id *iplayerDocument) programmeListSelection() *iplayerSelection {
	return &iplayerSelection{id.doc.Find(".content-item")}
}

// An IplayerDocumentResult is the result of generating a new
// goquery.Document from a iplayer url.
// If successful, the Idoc will be an iplayerDocument and nil for The error field.
type IplayerDocumentResult struct {
	Idoc  iplayerDocument
	Error error
}

type programPage struct {
	doc *iplayerDocument
}

func newProgramme(title, subtitle string, isel *iplayerSelection) *Programme {
	synopsis := isel.synopsis()
	thumbnail := isel.thumbnail()
	url := isel.url()
	available := isel.available()
	duration := isel.duration()
	return &Programme{
		title,
		subtitle,
		synopsis,
		thumbnail,
		url,
		0,
		available,
		duration,
	}
}

func (pp *programPage) programmes() []*Programme {
	var results []*Programme
	title := pp.doc.doc.Find(".hero-header__title").Text()
	pp.doc.doc.Find(".content-item").Each(func(i int, s *goquery.Selection) {
		subtitle := s.Find(".content-item__title").Text()
		results = append(results, newProgramme(title, subtitle, &iplayerSelection{s}))
	})
	return results
}

// A MainCategoryDocument is the collection point for an iplayer category.
// maindoc is the root (or page 1) document, nextdocs pages 2 - n,
// programpagedocs are the docuemnts for all programmes that have more
// than 1 episode, and selectionresults have the programmes with only
// one available episode.
type MainCategoryDocument struct {
	maindoc          *iplayerDocument
	nextdocs         []*iplayerDocument
	programpagedocs  []*iplayerDocument
	selectionresults []*iplayerSelectionResult
}

// Programmes traverses all iplayerdocuments of an MainCategoryDocument
// and returns all their programmes.
func (mcd *MainCategoryDocument) Programmes() []*Programme {
	var results []*Programme
	for _, i := range mcd.selectionresults {
		if i.prog != nil {
			results = append(results, i.prog)
		}
	}
	for _, i := range mcd.programpagedocs {
		pp := programPage{i}
		results = append(results, pp.programmes()...)
	}
	return results
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

func (id *iplayerDocument) programPages(selres []*iplayerSelectionResult) []Pager {
	var urls []Pager
	for _, i := range selres {
		if i.programPage != "" {
			urls = append(urls, BeebURL(i.programPage))
		}
	}
	return urls
}

type relatedLink struct {
	title string
	name  string
}

func (id *iplayerDocument) relatedLinks() []*relatedLink {
	var rellinks []*relatedLink
	id.doc.Find("related-link > a").Each(func(i int, s *goquery.Selection) {
		rl := relatedLink{s.Text(), s.AttrOr("href", "")}
		rellinks = append(rellinks, &rl)
	})
	return rellinks
}

// newMainCategory generates a new MainCategoryDocument
// from a root iplayer category document (eg. films, food)
func newMainCategory(np NextPager) *MainCategoryDocument {
	nextdocs := []*iplayerDocument{np.mainDoc()}
	var progpagedocs []*iplayerDocument
	npages := np.nextPages()
	nextPages := collectPages(npages)
	for _, i := range nextPages {
		if &i.Idoc != nil {
			nextdocs = append(nextdocs, &i.Idoc)
		} else {
			log.Fatal(&i.Error)
		}
	}
	var selres []*iplayerSelectionResult
	for _, i := range nextdocs {
		isel := i.programmeListSelection()
		selres = append(selres, isel.selectionResults()...)
	}
	urls := np.programPages(selres)
	progPages := collectPages(urls)
	for _, i := range progPages {
		if &i.Idoc != nil {
			progpagedocs = append(progpagedocs, &i.Idoc)
		} else {
			log.Fatal(&i.Error)
		}
	}
	return &MainCategoryDocument{np.mainDoc(), nextdocs[1:], progpagedocs, selres}
}

var seen = make(map[Pager]*IplayerDocumentResult)
var mutex = &sync.Mutex{}

func collectPages(urls []Pager) []*IplayerDocumentResult {
	var results []*IplayerDocumentResult
	c := make(chan *IplayerDocumentResult)
	jobs := 0
	for _, i := range urls {
		mutex.Lock()
		if res, ok := seen[i]; ok {
			mutex.Unlock()
			results = append(results, res)
		} else {
			mutex.Unlock()
			jobs++
			go func(u Pager) {
				u.loadDocument(c)
			}(i)
		}
	}
	for i := 0; i < jobs; i++ {
		res := <-c
		mutex.Lock()
		seen[res.Idoc.url] = res
		mutex.Unlock()
		results = append(results, res)
	}
	return results
}
