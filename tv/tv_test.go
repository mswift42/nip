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
	url = "testhtml/films1.html"
	go url.loadDocument(c)
	idr = <-c
	if idr.Error != nil {
		t.Error("Expected error to be nil: ", idr.Error)
	}
	if idr.Idoc.doc == nil {
		t.Error("Expected Idoc not to be nil: ", idr.Idoc)
	}
	url = "testhtml/nosuchfile.html"
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
	nsel := idr.Idoc.programmeNode()
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
		t.Error("Expected second programme Title to be 'Brideshead Revisited', got: ", selres[1].prog.Title)
	}
	if selres[1].programPage != "" {
		t.Error("Expected second programPage to be an empty string, got: ", selres[1].programPage)
	}
}

func TestCollectPages(t *testing.T) {
	doc := documentLoader("testhtml/films1.html")
	tid := TestIplayerDocument{doc}
	isel := tid.idoc.programmeNode()
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
	selres := tid.idoc.programmeNode().selectionResults()
	urls := tid.programPages(selres)
	if len(urls) != 2 {
		t.Error("Expected length of urls to be 2, got: ", len(urls))
	}
	if urls[0] != TestHTMLURL("testhtml/adam_curtis.html") {
		t.Error("Expected first Url to be 'adam_curtis', got: ", urls[0])
	}
	if urls[1] != TestHTMLURL("testhtml/storyville.html") {
		t.Error("Expected second Url to be 'storyville', got: ", urls[1])
	}
	doc = documentLoader("testhtml/food1.html")
	tid = TestIplayerDocument{doc}
	selres = tid.idoc.programmeNode().selectionResults()
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
	{
		"Series 1: 6. Sauces",
		"https://ichef.bbci.co.uk/images/ic/304x171/p062csj1.jpg",
		"Delia shares her tips on how to avoid lumps and curdling. (1978)",
		"/iplayer/episode/p05rt529/delia-smiths-cookery-course-series-1-6-sauces",
		"Available for over a year",
		"24 mins",
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
		t.Error("Expected Title of first storyville programme to be 'Storyville', "+
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
			t.Error("Expected Title to be Adam Curtis, got: ", i.Title)
		}
	}
	for i := range progs {
		if progs[i].Subtitle != AdamCurtis[i].subtitle {
			t.Errorf("Expected subtitle to be %q, got: %q", AdamCurtis[i].subtitle, progs[i].Subtitle)
		}
		if progs[i].URL != AdamCurtis[i].url {
			t.Errorf("Expected Url to be %q, got: %q", AdamCurtis[i].url, progs[i].URL)
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
			t.Errorf("Expected Title %q to have duration %q, got: %q",
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
			t.Error("Expected Title to be 'Delia Smith's Cookery Course, got: ",
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
			t.Errorf("Expecte Url to be %q, got: %q",
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

var filmurls = []string{
	"/iplayer/episode/b04n1hfy/storyville-112-weddings",
	"/iplayer/episode/p04b183c/adam-curtis-hypernormalisation",
	"/iplayer/episode/p02gyz6b/adam-curtis-bitter-lake",
	"/iplayer/episode/p0351g0z/fear-itself",
	"/iplayer/episode/b08nfjwt/wallace-and-gromit-the-wrong-trousers",
	"/iplayer/episode/b014f32p/brideshead-revisited",
	"/iplayer/episode/b00t61gx/the-damned-united",
	"/iplayer/episode/b0b49py6/dina",
	"/iplayer/episode/b007793l/great-expectations",
	"/iplayer/episode/b0b57pqy/kenny",
	"/iplayer/episode/b0b57d0w/king-lear",
	"/iplayer/episode/b0148wk1/ladies-in-lavender",
	"/iplayer/episode/b0074fln/lara-croft-tomb-raider",
	"/iplayer/episode/b007ck00/lara-croft-tomb-raider-the-cradle-of-life",
	"/iplayer/episode/b061tt33/lucky-them",
	"/iplayer/episode/b008lyr3/meet-the-fockers",
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
	"/iplayer/episode/p05rt529/delia-smiths-cookery-course-series-1-6-sauces",
	"/iplayer/episode/p05rt7cc/delia-smiths-cookery-course-series-1-7-spices-and-flavourings",
	"/iplayer/episode/p05rtnzc/delia-smiths-cookery-course-series-1-8-winter-vegetables",
	"/iplayer/episode/p05rtqvt/delia-smiths-cookery-course-series-1-9-pulses",
	"/iplayer/episode/p05rts0s/delia-smiths-cookery-course-series-1-10-puddings",
	"/iplayer/episode/b0b61r7j/britains-best-home-cook-series-1-episode-6",
	"/iplayer/episode/b0b5c1h5/britains-best-home-cook-series-1-episode-5",
	"/iplayer/episode/b0b4fkyb/britains-best-home-cook-series-1-episode-4",
	"/iplayer/episode/b0b3lh5g/britains-best-home-cook-series-1-episode-3",
	"/iplayer/episode/b0b2wj6k/britains-best-home-cook-series-1-episode-2",
	"/iplayer/episode/b0b2289t/britains-best-home-cook-series-1-episode-1",
	"/iplayer/episode/b0b53xqs/the-big-crash-diet-experiment",
	"/iplayer/episode/b08gj545/the-secrets-of-your-food-series-1-1-we-are-what-we-eat",
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
	{
		"Meet the Fockers",
		"",
		"/iplayer/episode/b008lyr3/meet-the-fockers",
		"Ben Stiller and Robert De Niro star in this blockbuster comedy sequel.",
		"https://ichef.bbci.co.uk/images/ic/304x171/p02x86xz.jpg",
		"Available until Sun 1:20am",
		"110 mins",
	},
	{
		"My Old Lady",
		"",
		"/iplayer/episode/b055d9vt/my-old-lady",
		"Comic, poignant relationship drama with Kevin Kline and Maggie Smith.",
		"https://ichef.bbci.co.uk/images/ic/304x171/p04tm316.jpg",
		"Available until Tue 12am",
		"102 mins",
	},
	{
		"Passion",
		"",
		"/iplayer/episode/b03ftl7s/passion",
		"Rivalry turns deadly in this tense thriller starring Rachel McAdams and Noomi Rapace.",
		"https://ichef.bbci.co.uk/images/ic/304x171/p03n5crh.jpg",
		"Expires tonight 2:35am",
		"94 mins",
	},
	{
		"The Past",
		"",
		"/iplayer/episode/b09sqss5/the-past",
		"A man attends his divorce hearing and gets embroiled in his wife's new relationship.",
		"https://ichef.bbci.co.uk/images/ic/304x171/p04rmr04.jpg",
		"Available for 17 days",
		"126 mins",
	},
	{
		"Quartet",
		"",
		"/iplayer/episode/b03ftm2k/quartet",
		"Ensemble comedy drama directed by Dustin Hoffman. With Maggie Smith and Tom Courtenay.",
		"https://ichef.bbci.co.uk/images/ic/304x171/p02fbt45.jpg",
		"Available for 24 days",
		"92 mins",
	},
	{
		"Shaun the Sheep the Movie",
		"",
		"/iplayer/episode/b05zxj5s/shaun-the-sheep-the-movie",
		"Shaun takes the day off to have some fun, but he gets more than he bargained for.",
		"https://ichef.bbci.co.uk/images/ic/304x171/p05r3g5y.jpg",
		"Available for 18 days",
		"85 mins",
	},
}

var foodprogs = []struct {
	title     string
	subtitle  string
	url       string
	synopsis  string
	thumbnail string
	available string
	duration  string
}{
	{
		"The Big Crash Diet Experiment",
		"",
		"/iplayer/episode/b0b53xqs/the-big-crash-diet-experiment",
		"Dr Javid Abdelmoneim and four overweight volunteers put crash dieting to the test.",
		"https://ichef.bbci.co.uk/images/ic/304x171/p067xj99.jpg",
		"Available for 21 days",
		"58 mins",
	},
	{
		"Britain's Best Home Cook",
		"Series 1: Episode 6",
		"/iplayer/episode/b0b61r7j/britains-best-home-cook-series-1-episode-6",
		"In the quarter-final the five remaining cooks must impress the judges with sharing feasts.",
		"https://ichef.bbci.co.uk/images/ic/304x171/p068kn5j.jpg",
		"Available for 6 months",
		"58 mins",
	},
	{
		"Britain's Fat Fight with Hugh Fearnley-Whittingstall",
		"Series 1: Episode 3",
		"/iplayer/episode/b0b2x6d8/britains-fat-fight-with-hugh-fearnleywhittingstall-series-1-episode-3",
		"Hugh learns of simple and obvious changes which could be made to GPs surgeries.",
		"https://ichef.bbci.co.uk/images/ic/304x171/p065x5ck.jpg",
		"Available for 3 months",
		"59 mins",
	},
	{
		"Britain's Fat Fight with Hugh Fearnley-Whittingstall",
		"Series 1: Episode 2",
		"/iplayer/episode/b0b1zh3y/britains-fat-fight-with-hugh-fearnleywhittingstall-series-1-episode-2",
		"Hugh Fearnley-Whittingstall turns the spotlight on popular high street restaurant chains.",
		"https://ichef.bbci.co.uk/images/ic/304x171/p0657lgn.jpg",
		"Available for 3 months",
		"59 mins",
	},
	{
		"Caribbean Food Made Easy",
		"Episode 4",
		"/iplayer/episode/b00mtr6m/caribbean-food-made-easy-episode-4",
		"Levi takes to the seas to show Mull fishermen how to dub up their catch.",
		"https://ichef.bbci.co.uk/images/ic/304x171/p05rgkdn.jpg",
		"Available for over a year",
		"29 mins",
	},
	{
		"Caribbean Food Made Easy",
		"Episode 3",
		"/iplayer/episode/b00mnk09/caribbean-food-made-easy-episode-3",
		"Levi challenges Cornish pasty eaters in Falmouth to try his golden vegetable patties.",
		"https://ichef.bbci.co.uk/images/ic/304x171/p05rgkgr.jpg",
		"Available for over a year",
		"29 mins",
	},
}

func TestNewMainCategory(t *testing.T) {
	doc := documentLoader("testhtml/food1.html")
	td := TestIplayerDocument{doc}
	nmc := newMainCategory(&td)
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
	foodprogrammes := nmc.Programmes()
	for _, i := range foodurls {
		if !contains(foodprogrammes, i) {
			t.Errorf("Expected foodprogrammes to contain %q ", i)
		}
	}
	for i := range foodprogs {
		found := findProgramme(foodprogrammes, foodprogs[i].url)
		if found == nil {
			t.Errorf("Expected programme: %q to be found.",
				foodprogs[i].title+foodprogs[i].subtitle)
		}
		if foodprogs[i].title != found.Title {
			t.Errorf("Expected programme to have have Title: %q. got: %q",
				foodprogs[i].title, found.Title)
		}
		if foodprogs[i].subtitle != found.Subtitle {
			t.Errorf("Expected programme %q to have subtitle: %q. got: %q ",
				foodprogs[i].title, foodprogs[i].subtitle, found.Subtitle)
		}
		if foodprogs[i].url != found.URL {
			t.Errorf("Expected programme %q %q to have Url: %q. got: %q ",
				foodprogs[i].title, foodprogs[i].subtitle, foodprogs[i].url, found.URL)
		}
		if foodprogs[i].thumbnail != found.Thumbnail {
			t.Errorf("Expected programme %q %q to have thumbnail: %q. got: %q ",
				foodprogs[i].title, foodprogs[i].subtitle, foodprogs[i].thumbnail, found.Thumbnail)
		}
		if foodprogs[i].available != found.Available {
			t.Errorf("Expected programme %q %q to have availability : %q. got: %q ",
				foodprogs[i].title, foodprogs[i].subtitle, foodprogs[i].available, found.Available)
		}
		if foodprogs[i].duration != found.Duration {
			t.Errorf("Expected programme %q %q to have duration %q, got: %q ",
				foodprogs[i].title, foodprogs[i].subtitle, foodprogs[i].duration, found.Duration)
		}
	}
	doc = documentLoader("testhtml/films1.html")
	td = TestIplayerDocument{doc}
	nmc = newMainCategory(&td)
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
			t.Errorf("Expected programme to have Title: %q. got: %q", filmprogs[i].title,
				found.Title)
		}
		if filmprogs[i].subtitle != found.Subtitle {
			t.Errorf("Expected programme to have subtitle: %q. got: %q",
				filmprogs[i].subtitle, found.Subtitle)
		}
		if filmprogs[i].url != found.URL {
			t.Errorf("Expected programme to have Url: %q. got: %q",
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
	cat := NewCategory("films", &td)
	if cat.Name != "films" {
		t.Errorf("Expected category's name to be 'films' , got: %q", cat.Name)
	}
	if len(cat.Programmes) != 28 {
		t.Error("Expected length of programmes to be 28, got: ", len(cat.Programmes))
	}
}

var tombRaiderProg = []struct {
	name string
	url  string
}{
	{
		"IMDb: Lara Croft Tomb Raider",
		"http://www.imdb.com/title/tt0146316/?ref_=nv_sr_1",
	},
	{
		"Rotten Tomatoes:  Lara Croft Tomb Raider",
		"https://www.rottentomatoes.com/m/lara_croft_tomb_raider",
	},
	{
		"Wikipedia: Lara Croft Tomb Raider",
		"https://en.wikipedia.org/wiki/Tomb_Raider",
	},
}

func TestRelatedLinks(t *testing.T) {
	doc := documentLoader("testhtml/lara_croft_tomb_raider_programme.html")
	td := TestIplayerDocument{doc}
	rl := td.idoc.relatedLinks()
	if len(rl) != 3 {
		t.Error("expected length of related links to be 3, got: ", len(rl))
	}
	for i := range tombRaiderProg {
		if tombRaiderProg[i].name != rl[i].Title {
			t.Errorf("Expected Title to be %q, got: %q",
				tombRaiderProg[i].name, rl[i].Title)
		}
		if tombRaiderProg[i].url != rl[i].URL {
			t.Errorf("Expected Url for %q to be %q, got: %q",
				tombRaiderProg[i].name, tombRaiderProg[i].url, rl[i].URL)
		}
	}
}

var catnametests = []struct {
	in  string
	out string
}{
	{"films", "films"},
	{"flms", "films"},
	{"food", "food"},
	{"comedy", "comedy"},
	{"cdy", "comedy"},
	{"crime", "crime"},
	{"clic", "classic+period"},
	{"classic", "classic+period"},
	{"scifi", "scifi+fantasy"},
	{"fantasy", "scifi+fantasy"},
	{"docu", "documentaries"},
	{"arts", "arts"},
	{"etainment", "entertainment"},
	{"hstory", "history"},
	{"life", "lifestyle"},
	{"music", "music"},
	{"news", "news"},
	{"science", "science+nature"},
	{"nature", "science+nature"},
	{"sport", "sport"},
	{"n", ""},
	{"", ""},
}

func TestCatNameCompleter(t *testing.T) {
	for _, i := range catnametests {
		out, _ := catNameCompleter(i.in)
		if out != i.out {
			t.Errorf("Expected out for %v to be %v, got %v",
				i.in, i.out, out)
		}
	}
}

func TestBoxSets(t *testing.T) {
	doc := documentLoader("testhtml/luther.html")
	td := TestIplayerDocument{doc}
	bs := td.idoc.boxSet()
	if len(bs) != 4 {
		t.Errorf("Expected length of boxSet to be 5, got: %v", len(bs))
	}
	if bs[0].String() != "/iplayer/episodes/b00vk2lp/luther?seriesId=b00vk2mq" {
		t.Error("expected, got ", bs[0].String())
	}
	if bs[1].String() != "/iplayer/episodes/b00vk2lp/luther?seriesId=p01b2b2g" {
		t.Errorf("Expected url to be %v, got: %v",
			"/iplayer/episodes/b00vk2lp/luther?seriesId=p01b2b2g",
			bs[1].String())
	}
	if bs[2].String() != "/iplayer/episodes/b00vk2lp/luther?seriesId=b06srp3h" {
		t.Errorf("expected url to be %v, got: %v",
			"/iplayer/episodes/b00vk2lp/luther?seriesId=b06srp3h",
			bs[2].String())
	}
	if bs[3].String() != "/iplayer/episodes/b00vk2lp/luther?seriesId=b0bxbh80" {
		t.Errorf("Expected url to be %v, got: %v",
			"/iplayer/episodes/b00vk2lp/luther?seriesId=b0bxbh80",
			bs[3].String())
	}
	doc = documentLoader("testhtml/fleabag.html")
	td = TestIplayerDocument{doc}
	bs = td.idoc.boxSet()
	if bs[0].String() != "/iplayer/episodes/p070npjv/fleabag?seriesId=p071bjr7" {
		t.Errorf("Expected url to be %v, got: %v",
			"/iplayer/episodes/p070npjv/fleabag?seriesId=p071bjr7",
			bs[0].String())
	}

	doc = documentLoader("testhtml/wrong_mans.html")
	td = TestIplayerDocument{doc}
	bs = td.idoc.boxSet()
	if bs[0].String() != "/iplayer/episodes/p02bhkmm/the-wrong-mans?seriesId=p02bhlq2" {
		t.Errorf("Expected url to be %v, got: %v",
			"/iplayer/episodes/p02bhkmm/the-wrong-mans?seriesId=p02bhlq2",
			bs[0].String())
	}
}
