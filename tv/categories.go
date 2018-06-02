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
