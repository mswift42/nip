package tv

import "fmt"

// Category represents an Iplayer Prograamme, consisting of its name
// and a list of its programmes.
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
	"drama+soaps":    BeebURL("https://www.bbc.co.uk/iplayer/categories/drama-and-soaps/a-z?sort=atoz"),
	"scifi+fantasy":  BeebURL("https://www.bbc.co.uk/iplayer/categories/drama-sci-fi-and-fantasy/a-z?sort=atoz"),
	"documentaries":  BeebURL("https://www.bbc.co.uk/iplayer/categories/documentaries/a-z?sort=atoz"),
	"arts":           BeebURL("https://www.bbc.co.uk/iplayer/categories/arts/a-z?sort=atoz"),
	"entertainment":  BeebURL("https://www.bbc.co.uk/iplayer/categories/entertainment/a-z?sort=atoz"),
	"history":        BeebURL("https://www.bbc.co.uk/iplayer/categories/history/a-z?sort=atoz"),
	"lifestyle":      BeebURL("https://www.bbc.co.uk/iplayer/categories/lifestyle/a-z?sort=atoz"),
	"music":          BeebURL("https://www.bbc.co.uk/iplayer/categories/music/a-z?sort=atoz"),
	"news":           BeebURL("https://www.bbc.co.uk/iplayer/categories/news/a-z?sort=atoz"),
	"science+nature": BeebURL("https://www.bbc.co.uk/iplayer/categories/science-and-nature/a-z?sort=atoz"),
	"sport":          BeebURL("https://www.bbc.co.uk/iplayer/categories/sport/a-z?sort=atoz"),
	"cbbc":           BeebURL("https://www.bbc.co.uk/iplayer/categories/cbbc/a-z?sort=atoz"),
}

func catNameCompleter(cat string) (string, error) {
	switch {
	case matchesName(cat, []string{"films", "flms", "film"}):
		return "films", nil
	case matchesName(cat, []string{"food", "fod", "fd"}):
		return "food", nil
	case matchesName(cat, []string{"comedy", "cdy", "come", "cmdy", "comed"}):
		return "comedy", nil
	case matchesName(cat, []string{"crime", "crm", "crim"}):
		return "crime", nil
	case matchesName(cat, []string{"classic+period", "classic", "period",
		"clsic", "prd", "class", "per", "peri", "clic"}):
		return "classic+period", nil
	case matchesName(cat, []string{"drama+soaps", "drama", "soaps", "dram",
		"drma", "dama", "soap", "sops", "sop"}):
		return "drama+soaps", nil
	case matchesName(cat, []string{"scifi+fantasy", "scifi", "fantasy", "scfi",
		"fanta", "sci", "scyfy", "scify", "fantas", "ftsy"}):
		return "scifi+fantasy", nil
	case matchesName(cat, []string{"documentaries", "docu", "documentary", "docus"}):
		return "documentaries", nil
	case matchesName(cat, []string{"arts", "art", "ats"}):
		return "arts", nil
	case matchesName(cat, []string{"entertainment", "etainment", "tainment",
		"etertainment"}):
		return "entertainment", nil
	case matchesName(cat, []string{"history", "hstory", "histo"}):
		return "history", nil
	case matchesName(cat, []string{"lifestyle", "lfestyle", "life", "style"}):
		return "lifestyle", nil
	case matchesName(cat, []string{"music", "msic", "music"}):
		return "music", nil
	case matchesName(cat, []string{"news", "nws"}):
		return "news", nil
	case matchesName(cat, []string{"science+nature", "science", "nature",
		"scence", "natur", "nture", "scienc", "sceince"}):
		return "science+nature", nil
	case matchesName(cat, []string{"sport", "spor", "sprt"}):
		return "sport", nil
	case matchesName(cat, []string{"cbbc", "cbb", "cbc", "ccbc", "ccb"}):
		return "cbbc", nil
	default:
		return "", fmt.Errorf("could not find any matching category")
	}
}

func matchesName(name string, nameslice []string) bool {
	for _, i := range nameslice {
		if name == i {
			return true
		}
	}
	return false
}

// NewCategory takes a category name and a category root document, generates
// a newMainCategory and returns a Category.
func NewCategory(name string, np NextPager) *Category {
	nmc := newMainCategory(np)
	return &Category{name, nmc.Programmes()}
}

func findCatTitle(url Pager) string {
	for k, v := range caturls {
		if url == v {
			return k
		}
	}
	return ""
}
