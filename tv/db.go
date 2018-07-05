package tv

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"fmt"

	"github.com/gosuri/uiprogress"
	"github.com/pkg/errors"
)

// ProgrammeDB represents a (file) DB of all saved
// Programmes, divided by Categories. The Saved field
// speciefies at what time the DB was last refreshed.
type ProgrammeDB struct {
	Categories []*Category `json:"categories"`
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

func (pdb *ProgrammeDB) findCategory(category string) (*Category, error) {
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
// Whenever a programmes Title contains the searchterm, the
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

// FindProgramme takes an index, queries the ProgrammeDB for it, and if found,
// returns it.
func (pdb *ProgrammeDB) FindProgramme(index int) (*Programme, error) {
	for _, i := range pdb.Categories {
		for _, j := range i.Programmes {
			if j.Index == index {
				return j, nil
			}
		}
	}
	return nil, fmt.Errorf("could not find Programme with index %d", index)
}

// FindURL takes an index, queries the ProgrammeDB for it, and if found,
// returns the URL for the matching programme.
func (pdb *ProgrammeDB) FindURL(index int) (string, error) {
	prog, err := pdb.FindProgramme(index)
	if err != nil {
		return "", err
	}
	return bbcprefix + prog.URL, nil
}

func (pdb *ProgrammeDB) FindRelatedLinks(index int) ([]*RelatedLink, error) {
	prog, err := pdb.FindProgramme(index)
	if err != nil {
		return nil, err
	}
	bu := BeebURL(prog.URL)
	c := make(chan *IplayerDocumentResult)
	go bu.loadDocument(c)
	idr := <-c
	if idr.Error != nil {
		return nil, err
	}
	hp := idr.Idoc.doc.Find(".inline-list__item > a").AttrOr("href", "")
	if hp == "" {
		return nil, fmt.Errorf("failed to find Programme Home Page")
	}
	bu = BeebURL(hp)
	go bu.loadDocument(c)
	idr = <-c
	if idr.Error != nil {
		return nil, err
	}
	return idr.Idoc.relatedLinks(), nil
}

// SaveDB makes a new Category for all entries in caturls,
// and if successful, stores stem in ProgrammeDB.
func SaveDB() {
	c := make(chan *IplayerDocumentResult)
	var np []NextPager
	var cats []*Category
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
