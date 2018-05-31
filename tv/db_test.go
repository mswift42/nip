package tv

import (
	"strings"
	"testing"
	"time"
)

func TestProgrammeDB_Index(t *testing.T) {
	doc := documentLoader("testhtml/films1.html")
	td := &TestIplayerDocument{doc}
	cat := NewCategory("films", td)
	pdb := &ProgrammeDB{[]*Category{cat}, time.Now()}
	pdb.index()
	fp := pdb.Categories[0].Programmes[0]
	if fp.Index != 0 {
		t.Errorf("Expected first programme to have index '0', got: %d", fp.Index)
	}
	for _, i := range pdb.Categories[0].Programmes[1:] {
		if !(i.Index > 0) {
			t.Errorf("Expected for title %q index to be '> 0', got: %d", i.Title, i.Index)
		}
	}

}

func TestProgrammeDB_Save(t *testing.T) {
	doc := documentLoader("testhtml/films1.html")
	td := &TestIplayerDocument{doc}
	cat := NewCategory("films", td)
	doc = documentLoader("testhtml/food1.html")
	td = &TestIplayerDocument{doc}
	cat2 := NewCategory("food", td)
	pdb := &ProgrammeDB{[]*Category{cat, cat2}, time.Now()}
	err := pdb.Save("mockdb.json")
	if err != nil {
		t.Error("Expected saving db should not return error, got: ", err)
	}
}

func TestRestoreProgrammeDB(t *testing.T) {
	pdb, err := RestoreProgrammeDB("mockdb.json")
	if err != nil {
		t.Errorf("Expected error to be nil, got: %q", err)
	}
	if pdb == nil {
		t.Error("Expected db not to be nil!")
	}
	if pdb.Categories[0].Name != "films" {
		t.Errorf("Expected first Category name to be 'films', got: %q ", pdb.Categories[0].Name)
	}
	if len(pdb.Categories) != 2 {
		t.Error("Expected length of categories to be 2, got: ", len(pdb.Categories))
	}
	if pdb.Categories[0].Programmes[0].Title != "A Simple Plan" {
		t.Errorf("Expected first programmes title to be 'A Simple Plan', got: %q ",
			pdb.Categories[0].Programmes[0].Title)
	}
}

func TestProgrammeDB_ListCategory(t *testing.T) {
	pdb, err := RestoreProgrammeDB("mockdb.json")
	if err != nil {
		t.Errorf("Expected error to be nil, got: %q ", err)
	}
	if pdb == nil {
		t.Error("Expected db not to be nil")
	}
	cat := pdb.ListCategory("films")
	if strings.Contains(cat, "Can't find Category") {
		t.Error("Expected ListCategory to find category films.")
	}
	if !strings.Contains(cat, "A Simple Plan") {
		t.Error("Expected ListCategory output to contain 'A Simple Plan'")
	}
	if !strings.Contains(cat, "Bill") {
		t.Error("Expected ListCategory output to contain 'Bill'")
	}
	nocat := pdb.ListCategory("flms")
	if !strings.Contains(nocat, "Can't find Category") {
		t.Error("Expected to get error message for missing Category.")
	}
	if strings.Contains(nocat, "A Simple Plan") {
		t.Error("There should be no listed Programmes.")
	}
	foodcat := pdb.ListCategory("food")
	if !strings.Contains(foodcat, "The Home That 2 Built") {
		t.Error("Expected ListCategory output to contain 'The Home That 2 Built'.")
	}
}

func TestProgrammeDB_FindTitle(t *testing.T) {
	pdb, err := RestoreProgrammeDB("mockdb.json")
	if err != nil {
		t.Errorf("Expected error to be nil, got: %q ", err)
	}
	prog := pdb.FindTitle("Bill")
	if !strings.Contains(prog, "Bill") {
		t.Error("Expected FindTitle to find Programme with title Bill.")
	}
	prog2 := pdb.FindTitle("Simple")
	if !strings.Contains(prog2, "A Simple Plan") {
		t.Error("Expected FindTitle to find Programme with title ' A Simple Plan '")
	}
	prog3 := pdb.FindTitle("The Home That 2 Built")
	lines := strings.Split(prog3, "\n")
	if len(lines) != 9 {
		t.Error("fxpected findTitle for 'The Home that 2 built', to be 9 lines, got: ",
			len(lines))
	}
	noprog := pdb.FindTitle("mnopqrst")
	if !strings.Contains(noprog, "No Matches found.\n") {
		t.Error("Did not expect to get a match.")
	}
}
