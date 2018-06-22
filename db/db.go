package db

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"fmt"
	"github.com/pkg/errors"
	"github.com/gosuri/uiprogress"
	"github.com/mswift42/nip/tv"
)

const bbcprefix = "https://bbc.co.uk"

// ProgrammeDB represents a (file) DB of all saved
// Programmes, divided by Categories. The Saved field
// speciefies at what time the DB was last refreshed.
type ProgrammeDB struct {
	Categories []*tv.Category `json:"categories"`
	Saved      time.Time   `json:"saved"`
}

// RestoreProgrammeDB takes a path to a json file, reads it, and if
// successful, unmarshals it as struct ProgrammeDB.
func RestoreProgrammeDB(filename string) (*ProgrammeDB, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var pdb ProgrammeDB
	json.Unmarshal(file, &pdb)
	return &pdb, nil
}

func (pdb *ProgrammeDB) toJSON() ([]byte, error) {
	marshalled, err := json.MarshalIndent(pdb, "", "\t")
	if err != nil {
		return nil, err
	}
	return marshalled, err
}

// Save takes a path to json file, adds the current time to field
// 'Saved', converts it to json in a human readable format, and if successful
// saves it to said file.
func (pdb *ProgrammeDB) Save(filename string) error {
	pdb.Saved = time.Now()
	pdb.index()
	enc, err := pdb.toJSON()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, enc, 0644)
}

func (pdb *ProgrammeDB) index() {
	index := 0
	for _, i := range pdb.Categories {
		for _, j := range i.Programmes {
			j.Index = index
			index++
		}
	}
}

// ListCategory takes the name for a category, searches the ProgrammeDB
// for it, and if found, returns a string with all the category's programmes.
func (pdb *ProgrammeDB) ListCategory(category string) string {
	var buffer bytes.Buffer
	cat, err := pdb.findCategory(category)
	if err != nil {
		return fmt.Sprintln(err)
	}
	for _, i := range cat.Programmes {
		buffer.WriteString(i.String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (pdb *ProgrammeDB) findCategory(category string) (*tv.Category, error) {
	for _, i := range pdb.Categories {
		if i.Name == category {
			return i, nil
		}
	}
	return nil, errors.New("Can't find Category with Name: " + category)
}

// ListAvailableCategories returns a string with the name of all categories stored
// in ProgrammeDB.
func (pdb *ProgrammeDB) ListAvailableCategories() string {
	var buffer bytes.Buffer
	for _, i := range pdb.Categories {
		buffer.WriteString(i.Name + "\n")
	}
	return buffer.String()
}
// FindTitle searches for a given term in the ProgrammeDB.
// Whenever a programmes title contains the searchterm, the
// programme is appended to the result string.
func (pdb *ProgrammeDB) FindTitle(title string) string {
	var buffer bytes.Buffer
	for _, i := range pdb.Categories {
		for _, j := range i.Programmes {
			if strings.Contains(strings.ToLower(j.String()),
				strings.ToLower(title)) {
				buffer.WriteString(j.String() + "\n")
			}
		}
	}
	if len(buffer.Bytes()) == 0 {
		return "No Matches found.\n"
	}
	return buffer.String()
}
func (pdb *ProgrammeDB) sixHoursLater(dt time.Time) bool {
	dur := dt.Sub(pdb.Saved)
	return dur.Truncate(time.Hour).Hours() >= 6
}

// FindURL takes an index, queries the ProgrammeDB for it, and if found,
// returns the URL for the matching programme.
func (pdb *ProgrammeDB) FindURL(index int) (string, error) {
	for _, i := range pdb.Categories {
		for _, j := range i.Programmes {
			if j.Index == index {
				return bbcprefix + j.URL, nil
			}
		}
	}
	return "", fmt.Errorf("could not find Programme with index %d", index)
}

// SaveDB makes a new Category for all entries in caturls,
// and if successful, stores stem in ProgrammeDB.
func SaveDB() {
	c := make(chan *tv.IplayerDocumentResult)
	var np []tv.NextPager
	var cats []*tv.Category
	uiprogress.Start()
	bar := uiprogress.AddBar(len(caturls)).AppendCompleted()
	for _, v := range caturls {
		go func(u Pager) {
			u.loadDocument(c)
		}(v)
	}
	for range caturls {
		docres := <-c
		if docres.Error == nil {
			np = append(np, &docres.Idoc)
			bar.Incr()
		} else {
			fmt.Println(docres.Error)
		}
	}
	for _, i := range np {
		nc := NewCategory(fincCatTitle(i.mainDoc().url), i)
		cats = append(cats, nc)
	}
	pdb := &ProgrammeDB{cats, time.Now()}
	uiprogress.Stop()
	pdb.Save("mockdb.json")
}

var caturls = map[string]tv.Pager{
	"films":          tv.BeebURL("https://www.bbc.co.uk/iplayer/categories/films/a-z?sort=atoz&page=1"),
	"food":           tv.BeebURL("https://www.bbc.co.uk/iplayer/categories/food/a-z?sort=atoz&page=1"),
	"comedy":         tv.BeebURL("https://www.bbc.co.uk/iplayer/categories/comedy/a-z?sort=atoz"),
	"crime":          tv.BeebURL("https://www.bbc.co.uk/iplayer/categories/drama-crime/a-z?sort=atoz"),
	"classic+period": tv.BeebURL("https://www.bbc.co.uk/iplayer/categories/drama-classic-and-period/a-z?sort=atoz"),
	"scifi+fantasy":  tv.BeebURL("https://www.bbc.co.uk/iplayer/categories/drama-sci-fi-and-fantasy/a-z?sort=atoz"),
	"documentaries":  tv.BeebURL("https://www.bbc.co.uk/iplayer/categories/documentaries/a-z?sort=atoz"),
	"arts":           tv.BeebURL("https://www.bbc.co.uk/iplayer/categories/arts/a-z?sort=atoz"),
	"entertainment":  tv.BeebURL("https://www.bbc.co.uk/iplayer/categories/entertainment/a-z?sort=atoz"),
	"history":        tv.BeebURL("https://www.bbc.co.uk/iplayer/categories/history/a-z?sort=atoz"),
	"lifestyle":      tv.BeebURL("https://www.bbc.co.uk/iplayer/categories/lifestyle/a-z?sort=atoz"),
	"music":          tv.BeebURL("https://www.bbc.co.uk/iplayer/categories/music/a-z?sort=atoz"),
	"news":           tv.BeebURL("https://www.bbc.co.uk/iplayer/categories/news/a-z?sort=atoz"),
	"science&nature": tv.BeebURL("https://www.bbc.co.uk/iplayer/categories/science-and-nature/a-z?sort=atoz"),
	"sport":          tv.BeebURL("https://www.bbc.co.uk/iplayer/categories/sport/a-z?sort=atoz"),
}
// NewCategory takes a category name and a category root document, generates
// a NewMainCategory and returns a Category.
func NewCategory(name string, np tv.NextPager) *tv.Category {
	nmc := tv.NewMainCategory(np)
	return &tv.Category{name, nmc.Programmes()}
}

func init() {
	pdb, err := RestoreProgrammeDB("mockdb.json")
	if err != nil {
		SaveDB()
	} else {
		if pdb.sixHoursLater(time.Now()) {
			SaveDB()
		}
	}
}
