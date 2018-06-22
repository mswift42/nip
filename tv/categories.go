package tv

// Category represents an Iplayer Prograamme, consisting of its name
// and a list of its programmes.
type Category struct {
	Name       string
	Programmes []*Programme
}

func loadCategory(name string, np NextPager, c chan<- *Category) {
	c <- NewCategory(name, np)
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
