package tv

import (
	"testing"
)

func TestLoadingDocument(t *testing.T) {
	url := TestHTMLURL("testhtml/food1.html")
	c := make(chan *IplayerDocumentResult)
	go url.loadDocument(c)
	idr := <-c
	if idr.Error != nil {
		t.Error("Expected error to be nil", idr.Error)
	}
	if idr.Idoc.doc == nil {
		t.Error("Expected Idoc not to be nil", idr.Idoc)
	}
	url = TestHTMLURL("testhtml/films1.html")
	go url.loadDocument(c)
	idr = <-c
	if idr.Error != nil {
		t.Error("Expected error to be nil: ", idr.Error)
	}
	if idr.Idoc.doc == nil {
		t.Error("Expected Idoc not to be nil: ", idr.Idoc)
	}
	url = TestHTMLURL("testhtml/nosuchfile.html")
	go url.loadDocument(c)
	idr = <-c
	if idr.Error == nil {
		t.Error("Expected to get error, got: ", idr.Error)
	}
}

func TestIplayerSelectionResults(t *testing.T) {
	url := TestHTMLURL("testhtml/films1.html")
	c := make(chan *IplayerDocumentResult)
	go url.loadDocument(c)
	idr := <-c
	sel := iplayerSelection{idr.Idoc.doc.Find(".content-item")}
	selres := sel.selectionResults()
	if len(selres) != 24 {
		t.Error("Expected length of selectionresults to equal: 24, got: ", len(selres))
	}
	nsel := idr.Idoc.programmeListSelection()
	nselres := nsel.selectionResults()
	if len(selres) != 24 {
		t.Error("Expected length of selectionResults to equal 24, got: ", len(nselres))
	}
	if selres[0].programPage != nselres[0].programPage {
		t.Error("Expected both selectionResults to be the same, got: ", nselres[0].programPage)
	}
	progpage := selres[0]
	if progpage.prog != nil {
		t.Error("Expected proramme to be nil: ", progpage.prog)
	}
	if progpage.programPage != "testhtml/adam_curtis.html" {
		t.Error("Expected program Page to be 'testhtml/adam_curtis.html' not: ", progpage.programPage)
	}
	if selres[1].prog.Title != "Brideshead Revisited" {
		t.Error("Expected second programme title to be 'Brideshead Revisited', got: ", selres[1].prog.Title)
	}
	if selres[1].programPage != "" {
		t.Error("Expected second programPage to be an empty string, got: ", selres[1].programPage)
	}
}

func TestCollectPages(t *testing.T) {
	doc := documentLoader("testhtml/films1.html")
	tid := TestIplayerDocument{doc}
	isel := tid.idoc.programmeListSelection()
	np := tid.programPages(isel.selectionResults())
	if len(np) != 2 {
		t.Error("Expected length of nextPages to be 1, got: ", len(np))
	}
	cp := collectPages(np)
	if len(cp) != 2 {
		t.Error("Expected length of collectedPages to be 1, got: ", len(cp))
	}
	if cp[0].Error != nil {
		t.Error("Expected error for first doc in collected Pages to be nil, got: ", cp[0].Error)
	}
}

func TestProgramPages(t *testing.T) {
	doc := documentLoader("testhtml/films1.html")
	tid := TestIplayerDocument{doc}
	selres := tid.idoc.programmeListSelection().selectionResults()
	urls := tid.programPages(selres)
	if len(urls) != 2 {
		t.Error("Expected length of urls to be 2, got: ", len(urls))
	}
	if urls[0] != TestHTMLURL("testhtml/adam_curtis.html") {
		t.Error("Expected first url to be 'adam_curtis', got: ", urls[0])
	}
	if urls[1] != TestHTMLURL("testhtml/storyville.html") {
		t.Error("Expected second url to be 'storyville', got: ", urls[1])
	}
	doc = documentLoader("testhtml/food1.html")
	tid = TestIplayerDocument{doc}
	selres = tid.idoc.programmeListSelection().selectionResults()
	urls = tid.programPages(selres)
	if len(urls) != 20 {
		t.Error("Expected length of urls to be 20, got: ", len(urls))
	}
	if urls[0] != TestHTMLURL("testhtml/britains_best_home_cook.html") {
		t.Error("Expected first food page to be 'britains_best_home_cook', got: ", urls[0])
	}
	if urls[19] != TestHTMLURL("testhtml/top_of_the_shop_with_tom_kerridge.html") {
		t.Error("expected last programpage to be 'top of the shop with "+
			"tom kerridge' got: ", urls[19])
	}
	docs := collectPages(urls)
	if len(docs) != 20 {
		t.Error("Expected length of collected docs to be 20, got: ", len(docs))
	}
	for _, i := range docs {
		if i.Error != nil {
			t.Error("Expected error to be nil, got: ", i.Error)
		}
	}
}

var AdamCurtis = []struct {
	subtitle  string
	thumbnail string
	synopsis  string
	url       string
	available string
	duration  string
}{
	{
		"HyperNormalisation",
		"https://ichef.bbci.co.uk/images/ic/304x171/p04c0tsb.jpg",
		"Welcome to the post-truth world. You know it’s not real. But you accept it as normal.",
		"/iplayer/episode/p04b183c/adam-curtis-hypernormalisation",
		"Available for over a year",
		"166 mins",
	},
	{
		"Bitter Lake",
		"https://ichef.bbci.co.uk/images/ic/304x171/p02h7n5x.jpg",
		"An adventurous and epic film by Adam Curtis.",
		"/iplayer/episode/p02gyz6b/adam-curtis-bitter-lake",
		"Available for over a year",
		"137 mins",
	},
}

// TODO Add food program page to test.
var DeliaSmith = []struct {
	subtitle  string
	thumbnail string
	synopsis  string
	url       string
	available string
	duration  string
}{
	{
		"Series 1: 10. Puddings",
		"https://ichef.bbci.co.uk/images/ic/304x171/p062csnk.jpg",
		"Delia makes several delicious and economical puddings. (1979)",
		"/iplayer/episode/p05rts0s/delia-smiths-cookery-course-series-1-10-puddings",
		"Available for over a year",
		"25 mins",
	},
	{
		"Series 1: 9. Pulses",
		"https://ichef.bbci.co.uk/images/ic/304x171/p062csn5.jpg",
		"Delia prepares unusual and economical meals from a range of pulses. (1979)",
		"/iplayer/episode/p05rtqvt/delia-smiths-cookery-course-series-1-9-pulses",
		"Available for over a year",
		"23 mins",
	},
	{
		"Series 1: 8. Winter Vegetables",
		"https://ichef.bbci.co.uk/images/ic/304x171/p062csjp.jpg",
		"Get the best out of winter veg with these delicious recipes. (1979)",
		"/iplayer/episode/p05rtnzc/delia-smiths-cookery-course-series-1-8-winter-vegetables",
		"Available for over a year",
		"24 mins",
	},
	{
		"Series 1: 7. Spices and Flavourings",
		"https://ichef.bbci.co.uk/images/ic/304x171/p062csj4.jpg",
		"Delia shows how to use spices to produce beautiful food. (1978)",
		"/iplayer/episode/p05rt7cc/delia-smiths-cookery-course-series-1-7-spices-and-flavourings",
		"Available for over a year",
		"25 mins",
	},
}

func TestProgramPage(t *testing.T) {
	doc := documentLoader("testhtml/storyville.html")
	pp := programPage{doc}
	progs := pp.programmes()
	if len(progs) != 4 {
		t.Error("Expected length of Programmes to be 4, got: ", len(progs))
	}
	if progs[0].Title != "Storyville" {
		t.Error("Expected title of first storyville programme to be 'Storyville', "+
			"got: ", progs[0].Title)
	}
	if progs[0].Subtitle != "112 Weddings" {
		t.Error("Expected subtitle of first storyville programme to be '112 Wedddings',"+
			"got: ", progs[0].Subtitle)
	}
	doc = documentLoader("testhtml/adam_curtis.html")
	pp = programPage{doc}
	progs = pp.programmes()
	if len(progs) != 2 {
		t.Error("Expected length of AdamCurtis programmes to be 2, got: ", len(progs))
	}
	for _, i := range progs {
		if i.Title != "Adam Curtis" {
			t.Error("Expected title to be Adam Curtis, got: ", i.Title)
		}
	}
	for i := range progs {
		if progs[i].Subtitle != AdamCurtis[i].subtitle {
			t.Errorf("Expected subtitle to be %q, got: %q", AdamCurtis[i].subtitle, progs[i].Subtitle)
		}
		if progs[i].URL != AdamCurtis[i].url {
			t.Errorf("Expected url to be %q, got: %q", AdamCurtis[i].url, progs[i].URL)
		}
		if progs[i].Thumbnail != AdamCurtis[i].thumbnail {
			t.Errorf("Expected thumbnail to be %q, got: %q", AdamCurtis[i].thumbnail, progs[i].Thumbnail)
		}
		if progs[i].Synopsis != AdamCurtis[i].synopsis {
			t.Errorf("Expected synopsis to be %q, got: %q", AdamCurtis[i].synopsis, progs[i].Synopsis)
		}
		if progs[i].Available != AdamCurtis[i].available {
			t.Errorf("Expected available to be %q, got: %q", AdamCurtis[i].available, progs[i].Available)
		}
		if progs[i].Duration != AdamCurtis[i].duration {
			t.Errorf("Expected title %q to have duration %q, got: %q",
				AdamCurtis[i].subtitle, AdamCurtis[i].duration, progs[i].Duration)
		}
	}
	doc = documentLoader("testhtml/delia_smiths_cookery_course.html")
	pp = programPage{doc}
	progs = pp.programmes()
	if len(progs) != 10 {
		t.Error("Expected length of Delia Smith programmes to be 10, got: ", len(progs))
	}
	for _, i := range progs {
		if i.Title != "Delia Smith's Cookery Course" {
			t.Error("Expected title to be 'Delia Smith's Cookery Course, got: ",
				i.Title)
		}
	}
	for i := range DeliaSmith {
		if progs[i].Subtitle != DeliaSmith[i].subtitle {
			t.Errorf("Epected subtitle to be %q, got: %q",
				DeliaSmith[i].subtitle, progs[i].Subtitle)
		}
		if progs[i].Synopsis != DeliaSmith[i].synopsis {
			t.Errorf("Expected synopsis to be %q, got: %q",
				DeliaSmith[i].synopsis, progs[i].Synopsis)
		}
		if progs[i].URL != DeliaSmith[i].url {
			t.Errorf("Expecte url to be %q, got: %q",
				DeliaSmith[i].url, progs[i].URL)
		}
		if progs[i].Thumbnail != DeliaSmith[i].thumbnail {
			t.Errorf("Expected thumbnail to be %q, got: %q",
				DeliaSmith[i].thumbnail, progs[i].Thumbnail)
		}
		if progs[i].Duration != DeliaSmith[i].duration {
			t.Errorf("Expected duration to be %q, got: %q",
				DeliaSmith[i].duration, progs[i].Duration)
		}
		if progs[i].Available != DeliaSmith[i].available {
			t.Errorf("Expected available to be %q, got: %q",
				DeliaSmith[i].available, progs[i].Available)
		}
	}
}

// TODO Add more urls to testurls.
var filmurls = []string{
	"/iplayer/episode/b04n1hfy/storyville-112-weddings",
	"/iplayer/episode/p04b183c/adam-curtis-hypernormalisation",
	"/iplayer/episode/p02gyz6b/adam-curtis-bitter-lake",
	"/iplayer/episode/p0351g0z/fear-itself",
	"/iplayer/episode/b08nfjwt/wallace-and-gromit-the-wrong-trousers",
}

var foodurls = []string{
	"/iplayer/episode/b00mtr6m/caribbean-food-made-easy-episode-4",
	"/iplayer/episode/p05rts0s/delia-smiths-cookery-course-series-1-10-puddings",
	"/iplayer/episode/p05rsy31/ken-homs-chinese-cookery-rice",
	"/iplayer/episode/p05rsw3r/ken-homs-chinese-cookery-meat",
	"/iplayer/episode/p05rr3hn/ken-homs-chinese-cookery-noodles",
	"/iplayer/episode/p05tjrrz/madhur-jaffreys-flavours-of-india-tamil-nadu",
	"/iplayer/episode/p05t9skh/madhur-jaffreys-flavours-of-india-goa",
	"/iplayer/episode/p05t9pn8/madhur-jaffreys-flavours-of-india-punjab",
	"/iplayer/episode/b01mwxk4/lorraines-fast-fresh-and-easy-food-6-everyday-easy",
	"/iplayer/episode/b01mrcxt/lorraines-fast-fresh-and-easy-food-5-posh-nosh",
	"/iplayer/episode/b01ml70w/lorraines-fast-fresh-and-easy-food-4-baking-it",
	"/iplayer/episode/b01mfxyy/lorraines-fast-fresh-and-easy-food-3-simple-classics",
}

// TODO add more progs to filmprogs/foodprogs.
var filmprogs = []struct {
	title     string
	subtitle  string
	url       string
	synopsis  string
	thumbnail string
	available string
	duration  string
}{
	{
		"Fear Itself",
		"",
		"/iplayer/episode/p0351g0z/fear-itself",
		"Uncover how films scare us with this mesmerising journey through horror cinema.",
		"https://ichef.bbci.co.uk/images/ic/304x171/p035db1t.jpg",
		"Available for over a year",
		"88 mins",
	},
	{
		"Lara Croft Tomb Raider: The Cradle of Life",
		"",
		"/iplayer/episode/b007ck00/lara-croft-tomb-raider-the-cradle-of-life",
		"Archaeologist Lara Croft faces a race against time to find mad bioweapons genius Dr Reiss.",
		"https://ichef.bbci.co.uk/images/ic/304x171/p05zxkj2.jpg",
		"Available for 2 months",
		"105 mins",
	},
}

func TestNewMainCategory(t *testing.T) {
	doc := documentLoader("testhtml/food1.html")
	td := TestIplayerDocument{doc}
	nmc := NewMainCategory(&td)
	if nmc.maindoc != td.idoc {
		t.Error("Expected maincategory maindoc and original document to be identical, got: ", nmc.maindoc)
	}
	if len(nmc.nextdocs) != 0 {
		t.Error("Expected length of nextdocs to be 1, got: ", len(nmc.nextdocs))
	}
	foodprogpagedocs := nmc.programpagedocs
	if len(foodprogpagedocs) != 20 {
		t.Error("Expected length of programpage docs to be 19, got: ", len(foodprogpagedocs))
	}
	foodprogs := nmc.Programmes()
	for _, i := range foodurls {
		if !contains(foodprogs, i) {
			t.Errorf("Expected foodprogs to contain %q ", i)
		}
	}
	doc = documentLoader("testhtml/films1.html")
	td = TestIplayerDocument{doc}
	nmc = NewMainCategory(&td)
	if len(nmc.nextdocs) != 0 {
		t.Error("Expected length of nextdocs to be 0, got: ", len(nmc.nextdocs))
	}
	filmprogpagedocs := nmc.programpagedocs
	if len(filmprogpagedocs) != 2 {
		t.Error("Expected length of film programpages to be 2, got: ", len(filmprogpagedocs))
	}
	filmres := nmc.selectionresults
	if len(filmres) != 24 {
		t.Error("Expected length of selectionresults to be 24, got: ", len(filmres))
	}
	filmprogrammes := nmc.Programmes()
	if len(filmprogrammes) != 28 {
		t.Error("Expected length of programmes to be 28, got: ", len(filmprogrammes))
	}
	for _, i := range filmurls {
		if !contains(filmprogrammes, i) {
			t.Errorf("Expected filmprogrammes to contain %q ", i)
		}
	}
	for i := range filmprogs {
		found := findProgramme(filmprogrammes, filmprogs[i].url)
		if found == nil {
			t.Errorf("Expected programme: %q to be found.",
				filmprogs[i].title+" "+filmprogs[i].subtitle)
		}
		if filmprogs[i].title != found.Title {
			t.Errorf("Expected programme to have title: %q. got: %q", filmprogs[i].title,
				found.Title)
		}
		if filmprogs[i].subtitle != found.Subtitle {
			t.Errorf("Expected programme to have subtitle: %q. got: %q",
				filmprogs[i].subtitle, found.Subtitle)
		}
		if filmprogs[i].url != found.URL {
			t.Errorf("Expected programme to have url: %q. got: %q",
				filmprogs[i].url, found.URL)
		}
		if filmprogs[i].synopsis != found.Synopsis {
			t.Errorf("Expected programme to have synopsis: %q. got: %q",
				filmprogs[i].synopsis, found.Synopsis)
		}
		if filmprogs[i].thumbnail != found.Thumbnail {
			t.Errorf("Expected programme to have thumbnail: %q. got: %q",
				filmprogs[i].thumbnail, found.Thumbnail)
		}
		if filmprogs[i].available != found.Available {
			t.Errorf("Expected programme: %q to have availability: %q. got: %q",
				filmprogs[i].title, filmprogs[i].available, found.Available)
		}
		if filmprogs[i].duration != found.Duration {
			t.Errorf("Expected programme: %q to have duration %q, got: %q",
				filmprogs[i].title, filmprogs[i].duration, found.Duration)
		}
	}
}

func TestCategory(t *testing.T) {
	doc := documentLoader("testhtml/films1.html")
	td := TestIplayerDocument{doc}
	cat := newCategory("films", &td)
	if cat.Name != "films" {
		t.Errorf("Expected category's name to be 'films' , got: %q", cat.Name)
	}
	if len(cat.Programmes) != 28 {
		t.Error("Expected length of programmes to be 28, got: ", len(cat.Programmes))
	}
}

func TestLoadCategories(t *testing.T) {
	doc := documentLoader("testhtml/films1.html")
	td := TestIplayerDocument{doc}
	doc2 := documentLoader("testhtml/food1.html")
	td2 := TestIplayerDocument{doc2}
	catmap := map[string]NextPager{
		"films": &td,
		"food":  &td2,
	}
	cats := loadCategories(catmap)
	if len(cats) != 2 {
		t.Error("Expected length of categories to be 2, got: ", len(cats))
	}
}
