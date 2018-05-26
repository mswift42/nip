package db

import (
	"github.com/mswift42/nip/tv"
	"time"
	"encoding/json"
	"io/ioutil"
)

type ProgrammeDB struct {
	Categories []tv.Category `json:"categories"`
	Saved time.Time `json:"saved"`
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