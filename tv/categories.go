package tv

type Category struct {
	Name       string
	Programmes []*Programme
}

var caturls = map[string]Pager{
	"films":          BeebURL("https://www.bbc.co.uk/iplayer/categories/films/a-z?sort=atoz&page=1"),
	"food":           BeebURL("https://www.bbc.co.uk/iplayer/categories/food/a-z?sort=atoz&page=1"),
	"comedy":         BeebURL("https://www.bbc.co.uk/iplayer/categories/comedy/a-z?sort=atoz"),
	"crime":          BeebURL("https://www.bbc.co.uk/iplayer/categories/drama-crime/a-z?sort=atoz"),
	"classic+period": BeebURL("https://www.bbc.co.uk/iplayer/categories/drama-classic-and-period/a-z?sort=atoz"),
	"scifi+fantasy":  BeebURL("https://www.bbc.co.uk/iplayer/categories/drama-sci-fi-and-fantasy/a-z?sort=atoz"),
	"documentaries":  BeebURL("https://www.bbc.co.uk/iplayer/categories/documentaries/a-z?sort=atoz"),
	"arts":           BeebURL("https://www.bbc.co.uk/iplayer/categories/arts/a-z?sort=atoz"),
	"entertainment":  BeebURL("https://www.bbc.co.uk/iplayer/categories/entertainment/a-z?sort=atoz"),
	"history":        BeebURL("https://www.bbc.co.uk/iplayer/categories/history/a-z?sort=atoz"),
	"lifestyle":      BeebURL("https://www.bbc.co.uk/iplayer/categories/lifestyle/a-z?sort=atoz"),
	"music":          BeebURL("https://www.bbc.co.uk/iplayer/categories/music/a-z?sort=atoz"),
	"news":           BeebURL("https://www.bbc.co.uk/iplayer/categories/news/a-z?sort=atoz"),
	"science&nature": BeebURL("https://www.bbc.co.uk/iplayer/categories/science-and-nature/a-z?sort=atoz"),
	"sport":          BeebURL("https://www.bbc.co.uk/iplayer/categories/sport/a-z?sort=atoz"),
}

func newCategory(name string, np NextPager) *Category {
	nmc := NewMainCategory(np)
	return &Category{name, nmc.Programmes()}
}

func loadCategory(name string, np NextPager, c chan<- *Category) {
	c <- newCategory(name, np)
}

func loadCategories(catmap map[string]NextPager) []*Category {
	var cats []*Category
	c := make(chan *Category)
	for n, np := range catmap {
		go func(name string, np NextPager) {
			loadCategory(name, np, c)
		}(n, np)
	}
	for range catmap {
		cats = append(cats, <-c)
	}
	return cats
}

func fincCatTitle(url Pager) string {
	for k, v := range caturls {
		if url == v {
			return k
		}
	}
	return ""
}
