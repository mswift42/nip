package tv

import (
	"fmt"
	"github.com/mswift42/goquery"
	"log"
	"regexp"
	"strings"
	"sync"
)

type BeebURL string

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
	//synopsis := is.synopsis()
	//url := is.url()
	//thumbnail := is.thumbNail()
	//pid := is.pid()
	//available := is.available()
	//duration := is.duration()
	//return &Programme{
	//	Title:     title,
	//	Subtitle:  subtitle,
	//	Synopsis:  synopsis,
	//	PID:       pid,
	//	Thumbnail: thumbnail,
	//	URL:       url,
	//	Index:     0,
	//	Available: available,
	//	Duration:  duration,
	//}
	return newProgrammeFromProgramPage(title, subtitle, is.sel)
}

func (is *iplayerSelection) title() string {
	return is.sel.Find(".content-item__title").Text()
}

func (is *iplayerSelection) subtitle() string {
	return is.sel.Find(".content-item__info__primary > " +
		".content-item__description").Text()
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

func (is *iplayerSelection) available() string {
	avail := is.sel.Find(".period").AttrOr("title", "")
	if avail == "" {
		return is.sel.Find(".availability-duration").Text()
	}
	return avail
}

func (is *iplayerSelection) duration() string {
	re := regexp.MustCompile(`\d+\smins`)
	return re.FindString(is.sel.Find(".duration").Last().Text())
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

type IplayerDocumentResult struct {
	Idoc  iplayerDocument
	Error error
}

type programPage struct {
	doc *iplayerDocument
}

func (pp *programPage) programmes() []*Programme {
	var results []*Programme
	title := pp.doc.doc.Find(".hero-header__title").Text()
	pp.doc.doc.Find(".content-item").Each(func(i int, s *goquery.Selection) {
		subtitle := s.Find(".content-item__title").Text()
		results = append(results, newProgrammeFromProgramPage(title, subtitle, s))
	})
	return results
}

func newProgrammeFromProgramPage(title string, subtitle string, s *goquery.Selection) *Programme {
	synopsis := s.Find(".content-item__info__secondary > .content-item__description").Text()
	url := s.Find("a").AttrOr("href", "")
	available := s.Find(".content-item__sublabels > span").Last().Text()
	duration := s.Find(".content-item__sublabels > span").First().Text()
	sel := iplayerSelection{s}
	thumbnail := sel.extractThumbnail()
	return &Programme{title, subtitle, synopsis, "",
		thumbnail, url, 0, available, duration}
}

func (is *iplayerSelection) extractThumbnail() string {
	set := is.sel.Find(".rs-image > picture > source").AttrOr("srcset", "")
	return strings.Split(set, " ")[0]
}

type MainCategoryDocument struct {
	maindoc          *iplayerDocument
	nextdocs         []*iplayerDocument
	programpagedocs  []*iplayerDocument
	selectionresults []*iplayerSelectionResult
}

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

//func DocumentsFromResults(docres []*IplayerDocumentResult) []*iplayerDocument {
//	var results []*iplayerDocument
//	for _, i := range docres {
//		if i.Error == nil {
//			results = append(results, &i.Idoc)
//		}
//	}
//	return results
//}

// TODO - update newMainCategory to use new iplayer maindoc layout.
func NewMainCategory(np NextPager) *MainCategoryDocument {
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
