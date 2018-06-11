package tv

import (
	"testing"
)
// TODO Update TESTHTML files.
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

//func TestIplayerSelectionResults(t *testing.T) {
//	url := TestHTMLURL("testhtml/films1.html")
//	c := make(chan *IplayerDocumentResult)
//	go url.loadDocument(c)
//	idr := <-c
//	sel := iplayerSelection{idr.Idoc.doc.Find(".list-item-inner")}
//	selres := sel.selectionResults()
//	if len(selres) != 20 {
//		t.Error("Expected length of selectionresults to equal: ", len(selres))
//	}
//	nsel := idr.Idoc.programmeListSelection()
//	nselres := nsel.selectionResults()
//	if len(selres) != 20 {
//		t.Error("Expected length of selectionResults to equal 20, got: ", len(nselres))
//	}
//	if selres[0].programPage != nselres[0].programPage {
//		t.Error("Expected both selectionResults to be the same, got: ", nselres[0].programPage)
//	}
//	progpage := selres[0]
//	if progpage.prog != nil {
//		t.Error("Expected proramme to be nil: ", progpage.prog)
//	}
//	if progpage.programPage != "testhtml/adam_curtis.html" {
//		t.Error("Expected program Page to be 'testhtml/adam_curtis.html' not: ", progpage.programPage)
//	}
//	if selres[1].prog.Title != "A Simple Plan" {
//		t.Error("Expected second programme title to be 'A Simple Plan', got: ", selres[1].prog.Title)
//	}
//	if selres[1].programPage != "" {
//		t.Error("Expected second programPage to be an empty string, got: ", selres[1].programPage)
//	}
//	url = TestHTMLURL("testhtml/films2.html")
//	go url.loadDocument(c)
//	idr = <-c
//	sel = iplayerSelection{idr.Idoc.doc.Find(".list-item-inner")}
//	selres = sel.selectionResults()
//	if len(selres) != 2 {
//		t.Error("Expected length of selectionresults to equal 2, got: ", len(selres))
//	}
//	if selres[0].prog.Title != "Wallace and Gromit: A Close Shave" {
//		t.Error("Expected title of first films2 programme to be wallace and gromit, got: ",
//			selres[0].prog.Title)
//	}
//	if selres[1].prog.Title != "Wonder Boys" {
//		t.Error("Expected title of second films2 programme to be 'Wonder Boys', got: ",
//			selres[1].prog.Title)
//	}
//}

//func TestCollectPages(t *testing.T) {
//	doc := documentLoader("testhtml/films1.html")
//	tid := TestIplayerDocument{doc}
//	np := tid.nextPages()
//	if len(np) != 1 {
//		t.Error("Expected length of nextPages to be 1, got: ", len(np))
//	}
//	cp := collectPages(np)
//	if len(cp) != 1 {
//		t.Error("Expected length of collectedPages to be 1, got: ", len(cp))
//	}
//	if cp[0].Error != nil {
//		t.Error("Expected error for first doc in collected Pages to be nil, got: ", cp[0].Error)
//	}
//	doc = documentLoader("testhtml/films2.html")
//	tid = TestIplayerDocument{doc}
//	np = tid.nextPages()
//	if len(np) != 1 {
//		t.Error("Expected length of nextPages to be 1, got: ", len(np))
//	}
//}

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
		t.Error("expected last programpage to be 'top of the shop with " +
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
}

var filmurls = []string{
	"/iplayer/episode/b04n1hfy/storyville-112-weddings",
	"/iplayer/episode/p04b183c/adam-curtis-hypernormalisation",
	"/iplayer/episode/p02gyz6b/adam-curtis-bitter-lake",
	"/iplayer/episode/b03p8shj/buena-vista-social-club",
	"/iplayer/episode/b00749zc/primary-colors",
	"/iplayer/episode/b0078nh3/wonder-boys",
	"/iplayer/episode/b0078cwc/a-simple-plan",
	"/iplayer/episode/p0351g0z/fear-itself",
	"/iplayer/episode/b05rmlr9/the-homesman",
	"/iplayer/episode/b01q0k5b/wallace-and-gromit-a-close-shave",
}

var foodurls = []string{
	"/iplayer/episode/b00mtr6m/caribbean-food-made-easy-episode-4",
	"/iplayer/episode/b0752bbd/chef-vs-science-the-ultimate-kitchen-challenge",
	"/iplayer/episode/b09yn368/classic-mary-berry-series-1-episode-6",
	"/iplayer/episode/p05rts0s/delia-smiths-cookery-course-series-1-10-puddings",
	"/iplayer/episode/p05rsy31/ken-homs-chinese-cookery-rice",
	"/iplayer/episode/p05rsw3r/ken-homs-chinese-cookery-meat",
	"/iplayer/episode/p05rr3hn/ken-homs-chinese-cookery-noodles",
	"/iplayer/episode/p05tjrrz/madhur-jaffreys-flavours-of-india-tamil-nadu",
	"/iplayer/episode/p05t9skh/madhur-jaffreys-flavours-of-india-goa",
	"/iplayer/episode/p05t9pn8/madhur-jaffreys-flavours-of-india-punjab",
	"/iplayer/episode/b07xsyr1/yes-chef-series-1-20-friday-final-4",
	"/iplayer/episode/b07xsq9v/yes-chef-series-1-19-mary-ann-gilchrist",
	"/iplayer/episode/b07xspv0/yes-chef-series-1-18-ryan-simpson",
	"/iplayer/episode/b07xsplr/yes-chef-series-1-17-atul-kochhar",
	"/iplayer/episode/b01mwxk4/lorraines-fast-fresh-and-easy-food-6-everyday-easy",
	"/iplayer/episode/b01mrcxt/lorraines-fast-fresh-and-easy-food-5-posh-nosh",
	"/iplayer/episode/b01ml70w/lorraines-fast-fresh-and-easy-food-4-baking-it",
	"/iplayer/episode/b01mfxyy/lorraines-fast-fresh-and-easy-food-3-simple-classics",
}

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
		"A Simple Plan",
		"",
		"/iplayer/episode/b0078cwc/a-simple-plan",
		"Bill Paxton gets caught up in lies, deceit and murder after the discovery of $4 million.",
		"https://ichef.bbci.co.uk/images/ic/336x189/p06586p5.jpg",
		"Available until 09:00 27 May 2018",
		"114 mins",
	},
	{
		"Bill",
		"",
		"/iplayer/episode/b08lvcg1/bill",
		"Bill Shakespeare (Matthew Baynton) leaves Stratford to follow his dream.",
		"https://ichef.bbci.co.uk/images/ic/336x189/p05r6x03.jpg",
		"Available until 19:00 23 May 2018",
		"88 mins",
	},
	{
		"Buena Vista Social Club",
		"",
		"/iplayer/episode/b03p8shj/buena-vista-social-club",
		"A group of Cuban musicians are brought together by Ry Cooder to record their music.",
		"https://ichef.bbci.co.uk/images/ic/336x189/p063zb3m.jpg",
		"Available until 22:40 20 May 2018",
		"105 mins",
	},
	{
		"Fear Itself",
		"",
		"/iplayer/episode/p0351g0z/fear-itself",
		"Uncover how films scare us with this mesmerising journey through horror cinema.",
		"https://ichef.bbci.co.uk/images/ic/336x189/p035db1t.jpg",
		"Available for over a year",
		"88 mins",
	},
	{
		"The Homesman",
		"",
		"/iplayer/episode/b05rmlr9/the-homesman",
		"Hilary Swank and Tommy Lee Jones star in this characterful western drama.",
		"https://ichef.bbci.co.uk/images/ic/336x189/p053038q.jpg",
		"Available until 00:50 16 May 2018",
		"113 mins",
	},
	{
		"Lara Croft Tomb Raider: The Cradle of Life",
		"",
		"/iplayer/episode/b007ck00/lara-croft-tomb-raider-the-cradle-of-life",
		"Archaeologist Lara Croft faces a race against time to find mad bioweapons genius Dr Reiss.",
		"https://ichef.bbci.co.uk/images/ic/336x189/p05zxkj2.jpg",
		"Available until 09:00 30 August 2018",
		"105 mins",
	},
	{
		"Man on the Moon",
		"",
		"/iplayer/episode/b007cjz1/man-on-the-moon",
		"Biopic of controversial comedian and star of Taxi and Saturday Night Live, Andy Kaufman.",
		"https://ichef.bbci.co.uk/images/ic/336x189/p05mwz8b.jpg",
		"Available until 09:00 27 May 2018",
		"118 mins",
	},
}

func TestNewMainCategory(t *testing.T) {
	doc := documentLoader("testhtml/food1.html")
	td := TestIplayerDocument{doc}
	nmc := NewMainCategory(&td)
	if nmc.maindoc != td.idoc {
		t.Error("Expected maincategory maindoc and original document to be identical, got: ", nmc.maindoc)
	}
	if len(nmc.nextdocs) != 1 {
		t.Error("Expected length of nextdocs to be 1, got: ", len(nmc.nextdocs))
	}
	food2 := nmc.nextdocs[0]
	isel := iplayerSelection{food2.doc.Find(".list-item-inner")}
	selres := isel.selectionResults()
	if selres[0].programPage != "testhtml/saturday_kitchen_best_bites.html" {
		t.Error("Expected 1st entry in food2 page to be 'Saturday Kitchen best bites', got: ",
			selres[0].programPage)
	}
	if len(selres) != 4 {
		t.Error("Expected length of selectionresults to be 4, got: ", len(selres))
	}
	for _, i := range selres {
		if i.prog != nil {
			t.Error("Expected prog to be nil, got: ", i.prog.Title)
		}
	}
	foodprogpagedocs := nmc.programpagedocs
	if len(foodprogpagedocs) != 19 {
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
	if len(nmc.nextdocs) != 1 {
		t.Error("Expected length of nextdocs to be 1, got: ", len(nmc.nextdocs))
	}
	filmprogpagedocs := nmc.programpagedocs
	if len(filmprogpagedocs) != 2 {
		t.Error("Expected length of film programpages to be 2, got: ", len(filmprogpagedocs))
	}
	filmres := nmc.selectionresults
	if len(filmres) != 22 {
		t.Error("Expected length of selectionresults to be 22, got: ", len(filmres))
	}
	filmprogrammes := nmc.Programmes()
	if len(filmprogrammes) != 26 {
		t.Error("Expected length of programmes to be 26, got: ", len(filmprogrammes))
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
	if len(cat.Programmes) != 26 {
		t.Error("Expected length of programmes to be 26, got: ", len(cat.Programmes))
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
