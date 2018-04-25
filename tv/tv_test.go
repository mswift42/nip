package tv

import (
	"testing"
)

func TestLoadingDocument(t *testing.T) {
	url := TestHTMLURL("testhtml/food1.html")
	c := make(chan *iplayerDocumentResult)
	go url.loadDocument(c)
	idr := <-c
	if idr.Error != nil {
		t.Error("Expected error to be nil", idr.Error)
	}
	if idr.idoc.doc == nil {
		t.Error("Expected idoc not to be nil", idr.idoc)
	}
	url = TestHTMLURL("testhtml/films1.html")
	go url.loadDocument(c)
	idr = <-c
	if idr.Error != nil {
		t.Error("Expected error to be nil: ", idr.Error)
	}
	if idr.idoc.doc == nil {
		t.Error("Expected idoc not to be nil: ", idr.idoc)
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
	c := make(chan *iplayerDocumentResult)
	go url.loadDocument(c)
	idr := <-c
	sel := iplayerSelection{idr.idoc.doc.Find(".list-item-inner")}
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
	if selres[1].prog.Title != "A Hijacking" {
		t.Error("Expected second programme title to be 'A Hijacking', got: ", selres[1].prog.Title)
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
		t.Error("Expected length of programmes to be 6, got: ", len(progs))
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
}
func TestNewMainCategory(t *testing.T) {
	doc := documentLoader("testhtml/films1.html")
	tid := TestIplayerDocument{doc}
	nmd := newMainCategory(&tid)
	if len(nmd.nextdocs) != 2 {
		t.Error("Expected length of nextdocs to be 2, got: ", len(nmd.nextdocs))
	}
	pp := programPage{nmd.nextdocs[1]}
	progs := pp.programmes()
	if len(progs) == 0 {
		t.Error("Expected length of programmes > 0, got: ", len(progs))
	}
	progs = nmd.programmes()
	if len(progs) != 28 {
		t.Error("Expected length of programmes for nmd to be 28, got: ", len(progs))
	}
	if !contains(progs, "/iplayer/episode/p04b183c/adam-curtis-hypernormalisation") {
		t.Error("Expected programmes to contain url of adam-curtis hypern.")
	}
	if !contains(progs, "/iplayer/episode/b041ycwk/a-hijacking") {
		t.Error("Expected programmes to contain url of A Hijacking.")
	}
	doc = documentLoader("testhtml/food1.html")
	tid  = TestIplayerDocument{doc}
	nmd = newMainCategory(&tid)
	if len(nmd.nextdocs) != 19 {
		t.Error("Expected length of nextdocs to be 19, got: ", len(nmd.nextdocs))
	}
	progs = nmd.programmes()
	if len(progs) != 88 {
		t.Error("Expected length of programmes to be 88, got: ", len(progs))
	}
}


