package tv

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"fmt"
	"github.com/pkg/errors"
)

// ProgrammeDB represents a (file) DB of all saved
// Programmes, divided by Categories. The Saved field
// speciefies at what time the DB was last refreshed.
type ProgrammeDB struct {
	Categories []*Category `json:"categories"`
	Saved      time.Time   `json:"saved"`
}

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

func (pdb *ProgrammeDB) ListAvailableCategories() string {
	var buffer bytes.Buffer
	for _, i := range pdb.Categories {
		buffer.WriteString(i.Name + "\n")
	}
	return buffer.String()
}

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

func (pdb *ProgrammeDB) sixHoursLater() bool {
	return time.Since(pdb.Saved).Minutes() < 6 * 60
}
