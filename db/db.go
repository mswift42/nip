package db

import (
	"github.com/mswift42/nip/tv"
	"time"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
)

type ProgrammeDB struct {
	Categories []tv.Category `json:"categories"`
	Saved time.Time `json:"saved"`
}

func RestoreProgrammeDB(filename string) (*ProgrammeDB, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var pdb *ProgrammeDB
	json.Unmarshal(file, pdb)
	return pdb, nil
}

func (pdb *ProgrammeDB) toJson() ([]byte, error) {
	marshalled, err := json.MarshalIndent(pdb, "", "\t")
	if err != nil {
		return nil, err
	}
	return marshalled, err
}

func (pdb *ProgrammeDB) Save(filename string) error {
	pdb.Saved = time.Now()
	pdb.index()
	json ,err := pdb.toJson()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, json, 0644)
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

func (pdb *ProgrammeDB) findCategory(category string) (*tv.Category, error) {
	for _, i := range pdb.Categories {
		if i.Name ==  category {
			return i, nil
		}
	}
	return nil, errors.New("Can't find Category with Name: " + category)
}