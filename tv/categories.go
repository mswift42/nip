package tv

type Category struct {
	Name       string
	Programmes []*Programme
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
