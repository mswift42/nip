package tv

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"fmt"

	"log"

	"github.com/gosuri/uiprogress"
)

// ProgrammeDB represents a (file) DB of all saved
// Programmes, divided by Categories. The Saved field
// speciefies at what time the DB was last refreshed.
type ProgrammeDB struct {
	Categories      []*Category       `json:"categories"`
	Saved           time.Time         `json:"saved"`
	SavedProgrammes []*SavedProgramme `json:"saved_programmes"`
}

// SavedProgramme is the url to a downloaded programme
// and the time it was downloaded.
type SavedProgramme struct {
	File  string    `json:"url"`
	Saved time.Time `json:"saved"`
}

// RestoreProgrammeDB takes a path to a json file, reads it, and if
// successful, unmarshals it as struct ProgrammeDB.
func RestoreProgrammeDB(filename string) (*ProgrammeDB, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var pdb ProgrammeDB
	err = json.Unmarshal(file, &pdb)
	if err != nil {
		return nil, err
	}
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
	cat, err := catNameCompleter(category)
	if err == nil {
		for _, i := range pdb.Categories {
			if cat == i.Name {
				return i, nil
			}
		}
	}
	return nil, err
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

func (pdb *ProgrammeDB) toBeDeletedProgrammes() []*SavedProgramme {
	var sp []*SavedProgramme
	for _, i := range pdb.SavedProgrammes {
		fmt.Println("Ranging over Saved Programmes: ", i)
		if _, err := os.Stat(i.File); os.IsExist(err) {
			since := time.Since(i.Saved).Truncate(time.Hour).Hours() / 24
			if since > 30.0 {
				sp = append(sp, i)
			}
		}
	}
	return sp
}

func (pdb *ProgrammeDB) removeFromSaved(sp *SavedProgramme) {
	var progs []*SavedProgramme
	for _, i := range pdb.SavedProgrammes {
		if sp.File != i.File {
			progs = append(progs, i)
		}
	}
	pdb.SavedProgrammes = progs
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
	return BBCPrefix + prog.URL, nil
}

// MarkSaved adds the filename of a downloaded programme + the
// date when it was downloaded to the SavedProgrammes entry in the ProgrammeDB.
func (pdb *ProgrammeDB) MarkSaved(filename string) {
	path := GetDBPath()
	pdbold, err := RestoreProgrammeDB(path + NipDB)
	if err != nil {
		log.Fatal(err)
	}
	sp := append(pdbold.SavedProgrammes, &SavedProgramme{filename, time.Now()})
	pdbnew := &ProgrammeDB{pdbold.Categories, pdb.Saved, sp}
	pdbnew.Save(path + NipDB)
}

// FindRelatedLinks loads the root page of a Programme.
// If found, it returns a slice of its related links, e.g. IMDB, Wikipedia, RottenTomatoes,...
func (pdb *ProgrammeDB) FindRelatedLinks(index int) ([]*RelatedLink, error) {
	prog, err := pdb.FindProgramme(index)
	if err != nil {
		return nil, err
	}
	uiprogress.Start()
	bar := uiprogress.AddBar(2).AppendCompleted()
	bu := BeebURL(prog.URL)
	c := make(chan *IplayerDocumentResult)
	go bu.loadDocument(c)
	idr := <-c
	bar.Incr()
	if idr.Error != nil {
		return nil, err
	}
	hp := idr.Idoc.doc.Find(".inline-list__item > a").AttrOr("href", "")
	if hp == "" {
		return nil, fmt.Errorf("failed to find Programme Home Page")
	}
	bu = BeebURL(BBCPrefix + hp)
	go bu.loadDocument(c)
	idr = <-c
	bar.Incr()
	if idr.Error != nil {
		return nil, err
	}
	uiprogress.Stop()
	return idr.Idoc.relatedLinks(), nil
}

// RefreshDB makes a new Category for all entries in caturls,
// and if successful, stores stem in ProgrammeDB.
func RefreshDB(filename string) {
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
		nc := NewCategory(findCatTitle(i.mainDoc().url), i)
		cats = append(cats, nc)
	}
	pdbold, err := RestoreProgrammeDB(filename)
	if err != nil {
		panic(err)
	}
	sp := pdbold.SavedProgrammes
	pdb := &ProgrammeDB{cats, time.Now(), sp}
	uiprogress.Stop()
	err = pdb.Save(filename)
	if err != nil {
		panic(err)
	}
}

func init() {
	dbpath := GetDBPath()
	filename := NipDB
	pdb, err := RestoreProgrammeDB(dbpath + filename)
	if err != nil {
		RefreshDB(dbpath + filename)
	} else {
		if pdb.sixHoursLater(time.Now()) || len(pdb.Categories) == 0 {
			RefreshDB(dbpath + filename)
		}
	}
	tobedel := pdb.toBeDeletedProgrammes()
	if len(tobedel) > 0 {
		fmt.Println("The following Programmes were downloaded > 30 days ago and have to be deleted: ")
		for _, i := range tobedel {
			fmt.Println(i.File)
		}
	}
}
