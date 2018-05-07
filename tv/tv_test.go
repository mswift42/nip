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
	sel := iplayerSelection{idr.Idoc.doc.Find(".list-item-inner")}
	selres := sel.selectionResults()
	if len(selres) != 20 {
		t.Error("Expected length of selectionresults to equal: ", len(selres))
	}
	progpage := selres[0]
	if progpage.prog != nil {
		t.Error("Expected proramme to be nil: ", progpage.prog)
	}
	if progpage.programPage != "testhtml/adam_curtis.html" {
		t.Error("Expected program Page to be 'testhtml/adam_curtis.html' not: ", progpage.programPage)
	}
	if selres[1].prog.Title != "A Simple Plan" {
		t.Error("Expected second programme title to be 'A Simple Plan', got: ", selres[1].prog.Title)
	}
	if selres[1].programPage != "" {
		t.Error("Expected second programPage to be an empty string, got: ", selres[1].programPage)
	}
}

func TestCollectPages(t *testing.T) {
	doc := documentLoader("testhtml/films1.html")
	tid := TestIplayerDocument{doc}
	np := tid.nextPages()
	if len(np) != 1 {
		t.Error("Expected length of nextPages to be 1, got: ", len(np))
	}
	cp := collectPages(np)
	if len(cp) != 1 {
		t.Error("Expected length of collectedPages to be 1, got: ", len(cp))
	}
	if cp[0].Error != nil {
		t.Error("Expected error for first doc in collected Pages to be nil, got: ", cp[0].Error)
	}
}

var classicMary = []struct {
	subtitle  string
	thumbnail string
	synopsis  string
	url       string
}{
	{
		"Series 1: Episode 6",
		"https://ichef.bbci.co.uk/images/ic/304x171/p062dlmz.jpg",
		"Mary unleashes some of her classic favourites that have made a comeback.",
		"/iplayer/episode/b09yn368/classic-mary-berry-series-1-episode-6",
	},
	{
		"Series 1: Episode 5",
		"https://ichef.bbci.co.uk/images/ic/304x171/p061mhz1.jpg",
		"Mary Berry returns to one of her most loved locations - Port Isaac in Cornwall.",
		"/iplayer/episode/b09xsw6b/classic-mary-berry-series-1-episode-5",
	},
	{
		"Series 1: Episode 4",
		"https://ichef.bbci.co.uk/images/ic/304x171/p06106t8.jpg",
		"Mary has always loved entertaining and creates some timeless classics.",
		"/iplayer/episode/b09x0tfw/classic-mary-berry-series-1-episode-4",
	},
	{
		"Series 1: Episode 3",
		"https://ichef.bbci.co.uk/images/ic/304x171/p06084xr.jpg",
		"Mary embraces the countryside with cooking inspired by food grown on farms and in gardens.",
		"/iplayer/episode/b09w3ynk/classic-mary-berry-series-1-episode-3",
	},
	{
		"Series 1: Episode 2",
		"https://ichef.bbci.co.uk/images/ic/304x171/p05zf2vg.jpg",
		"Mary Berry takes inspiration from a visit to a groundbreaking primary school in London.",
		"/iplayer/episode/b09vfd5d/ad/classic-mary-berry-series-1-episode-2",
	},
	{
		"Series 1: Episode 1",
		"https://ichef.bbci.co.uk/images/ic/304x171/p05yp3kv.jpg",
		"Mary Berry indulges her love of comfort food with homely recipes.",
		"/iplayer/episode/b09tp4ff/ad/classic-mary-berry-series-1-episode-1",
	},
}

func TestProgramPage(t *testing.T) {
	doc := documentLoader("testhtml/classic_mary_berry.html")
	pp := programPage{doc}
	progs := pp.programmes()
	if len(progs) != 6 {
		t.Error("Expected length of Programmes to be 6, got: ", len(progs))
	}
	for _, i := range progs {
		if i.Title != "Classic Mary Berry" {
			t.Error("Expected Title to be 'Classic Mary Berry, got: ", i.Title)
		}
	}
	for i := range classicMary {
		if progs[i].Subtitle != classicMary[i].subtitle {
			t.Error("Expected subtitle to be : "+classicMary[i].subtitle+" got: ", progs[i].Subtitle)
		}
		if progs[i].Synopsis != classicMary[i].synopsis {
			t.Error("Expected synopsis to be : "+classicMary[i].synopsis+" gog: ", progs[i].Synopsis)
		}
		if progs[i].URL != classicMary[i].url {
			t.Error("Expected url to be: "+classicMary[i].url+" got: ", progs[i].URL)
		}
	}
	doc = documentLoader("testhtml/storyville.html")
	pp = programPage{doc}
	progs = pp.programmes()
	if len(progs) != 4 {
		t.Error("Expected length of Programmes to be 4, got: ", len(progs))
	}
	if progs[0].Title != "Storyville" {
		t.Error("Expected title of first storyville programme to be 'Storyville', "+
			"got: ", progs[0].Title)
	}
	if progs[0].Subtitle != "112 Weddings" {
		t.Error("Expected subtitle of first storyville programme to be '112 Wedddings'," +
			"got: ", progs[0].Subtitle)
	}
}

var filmurls = []string {
		"/iplayer/episode/b04n1hfy/storyville-112-weddings",
		"/iplayer/episode/p04b183c/adam-curtis-hypernormalisation",
		"/iplayer/episode/p02gyz6b/adam-curtis-bitter-lake",
		"/iplayer/episode/b01q0k5b/wallace-and-gromit-a-close-shave",
		"/iplayer/episode/b08nfjwt/wallace-and-gromit-the-wrong-trousers",
}

var foodurls = []string {
		"/iplayer/episode/b00mtr6m/caribbean-food-made-easy-episode-4",
		"/iplayer/episode/b0752bbd/chef-vs-science-the-ultimate-kitchen-challenge",
		"/iplayer/episode/b09yn368/classic-mary-berry-series-1-episode-6",
		"/iplayer/episode/p05rts0s/delia-smiths-cookery-course-series-1-10-puddings",
}

func TestNewMainCategory(t *testing.T) {
	doc := documentLoader("testhtml/films1.html")
	tid := TestIplayerDocument{doc}
	np := tid.nextPages()
	if len(np) != 1 {
		t.Error("Expected length of nextpages to be 1, got: ", len(np))
	}
	if np[0] != TestHTMLURL("testhtml/films2.html") {
		t.Error("Expected nextpage to be 'testhtml/films2.html', got: ", np[0])
	}
	nmd := NewMainCategory(&tid)
	if len(nmd.programpagedocs) != 2 {
		t.Error("Expected length of programpagedocs to be 2, got: ", len(nmd.programpagedocs))
	}
	progs := nmd.Programmes()
	if len(progs) != 21 {
		t.Error("Expected length of film Programmes to be 27, got: ", len(progs))
	}
	for _, i := range filmurls {
		if !contains(progs, i) {
			t.Errorf("Expected filmurls to contain %s ", i)
		}
	}
	doc = documentLoader("testhtml/food1.html")
	tid = TestIplayerDocument{doc}
	nmd = NewMainCategory(&tid)
	if len(nmd.nextdocs) != 1 {
		t.Error("Expected length of nextdocs to be 19, got: ", len(nmd.nextdocs))
	}
	if len(nmd.programpagedocs) != 19 {
		t.Error("Expected length of programPage docs to be 19, got: ", len(nmd.programpagedocs))
	}
	if len(nmd.selectionresults) != 24 {
		t.Error("Expected length of selectionresults to be 24, got: ", len(nmd.selectionresults))
	}
}
