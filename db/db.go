package db

import (
	"github.com/mswift42/nip/tv"
	"time"
)

type ProgrammeDB struct {
	Categories []tv.Category `json:"categories"`
	Saved time.Time `json:"saved"`
}