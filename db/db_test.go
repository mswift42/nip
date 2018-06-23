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
	npdb, err := RestoreProgrammeDB("nosuchfile.json")
	if npdb != nil {
		t.Error("Expected db to be nil, got: ", npdb)
	}
	if err == nil {
		t.Error("Expected err not to be nil, got: ", err)
	}
	if pdb.Categories[0].Name != "films" {
		t.Errorf("Expected first Category name to be 'films', got: %q ", pdb.Categories[0].Name)
	}
	if len(pdb.Categories) != 2 {
		t.Error("Expected length of categories to be 2, got: ", len(pdb.Categories))
	}
	if pdb.Categories[0].Programmes[0].Title != "Brideshead Revisited" {
		t.Errorf("Expected first programmes title to be 'Brideshead Revisited', got: %q ",
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
	if !strings.Contains(cat, "Brideshead Revisited") {
		t.Error("Expected ListCategory output to contain 'Brideshead Revisited'")
	}
	if !strings.Contains(cat, "Tomb Raider") {
		t.Error("Expected ListCategory output to contain 'Tomb Raider'")
	}
	nocat := pdb.ListCategory("flms")
	if !strings.Contains(nocat, "Can't find Category") {
		t.Error("Expected to get error message for missing Category.")
	}
	if strings.Contains(nocat, "Brideshead") {
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
	prog := pdb.FindTitle("Tomb")
	if !strings.Contains(prog, "Tomb") {
		t.Error("Expected FindTitle to find Programme with title Bill.")
	}
	prog2 := pdb.FindTitle("Brideshead")
	if !strings.Contains(prog2, "Brideshead") {
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

func TestProgrammeDB_SixHoursLater(t *testing.T) {
	pdb, err := RestoreProgrammeDB("mockdb.json")
	if err != nil {
		t.Errorf("Expected error to be nil, got|: %s ", err)
	}
	t1 := time.Date(2018, 6, 16, 12, 0, 0, 0, time.UTC)
	t2 := time.Date(2018, 6, 16, 18, 0, 0, 0, time.UTC)
	t3 := time.Date(2018, 6, 16, 17, 59, 0, 0, time.UTC)
	t4 := time.Date(2018, 6, 18, 0, 0, 0, 0, time.UTC)
	t5 := time.Date(2018, 6, 16, 13, 0, 0, 0, time.UTC)
	pdb.Saved = t1
	if !pdb.sixHoursLater(t2) {
		t.Error("Expected sixHoursLater to be true, got: ", pdb.sixHoursLater(t2))
	}
	if pdb.sixHoursLater(t3) {
		t.Error("Expected sixHoursLater to be false, got: ", pdb.sixHoursLater(t3))
	}
	if !pdb.sixHoursLater(t4) {
		t.Error("Expected sixHoursLater to be true, got: ", pdb.sixHoursLater(t4))
	}
	if pdb.sixHoursLater(t5) {
		t.Error("Expected sixHoursLater to be false, got: ", pdb.sixHoursLater(t5))
	}

}
