package tv

import (
	"testing"
	"time"
)

func TestProgrammeDB_Save(t *testing.T) {
	doc := documentLoader("testhtml/films1.html")
	td := &TestIplayerDocument{doc}
	cat := NewCategory("films", td)
	pdb := &ProgrammeDB{[]*Category{cat}, time.Now()}
	err := pdb.Save("testdb.json")
	if err != nil {
		t.Error("Expected saving db should not return error, got: ", err)
	}
}
