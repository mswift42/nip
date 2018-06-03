package tv

type Category struct {
	Name       string
	Programmes []*Programme
}

var caturls = map[string]string{
	"films" : "https://www.bbc.co.uk/iplayer/categories/films/a-z?sort=atoz&page=1",
	"food" : "https://www.bbc.co.uk/iplayer/categories/food/a-z?sort=atoz&page=1",
	"comedy" : "https://www.bbc.co.uk/iplayer/categories/comedy/a-z?sort=atoz",
	"crime" : "https://www.bbc.co.uk/iplayer/categories/drama-crime/a-z?sort=atoz",
	"classic&period" : "https://www.bbc.co.uk/iplayer/categories/drama-classic-and-period/a-z?sort=atoz",
	"sc-fi&fantasy" : "https://www.bbc.co.uk/iplayer/categories/drama-sci-fi-and-fantasy/a-z?sort=atoz",
	"documentaries" : "https://www.bbc.co.uk/iplayer/categories/documentaries/a-z?sort=atoz",
	"arts" : "https://www.bbc.co.uk/iplayer/categories/arts/a-z?sort=atoz",
	"entertainment" : "https://www.bbc.co.uk/iplayer/categories/entertainment/a-z?sort=atoz",
	"history" : "https://www.bbc.co.uk/iplayer/categories/history/a-z?sort=atoz",
	"lifestyle" : "https://www.bbc.co.uk/iplayer/categories/lifestyle/a-z?sort=atoz",
	"music" : "https://www.bbc.co.uk/iplayer/categories/music/a-z?sort=atoz",
	"news"  : "https://www.bbc.co.uk/iplayer/categories/news/a-z?sort=atoz",
	"science&nature" : "https://www.bbc.co.uk/iplayer/categories/science-and-nature/a-z?sort=atoz",
	"sport" : "https://www.bbc.co.uk/iplayer/categories/sport/a-z?sort=atoz",
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
