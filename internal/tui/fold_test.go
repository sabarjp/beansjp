package tui

import "testing"

func newFoldTestModel() listModel {
	return listModel{
		foldOverrides: map[string]bool{},
		autoFolded:    map[string]bool{"e2": true},
	}
}

func TestIsFolded(t *testing.T) {
	m := newFoldTestModel()

	if !m.isFolded("e2") {
		t.Error("e2: expected auto-folded")
	}
	if m.isFolded("e1") {
		t.Error("e1: expected unfolded")
	}

	// Explicit toggles win over the automatic state, in both directions
	m.foldOverrides["e2"] = false
	m.foldOverrides["e1"] = true

	if m.isFolded("e2") {
		t.Error("e2: explicit unfold should win over auto-fold")
	}
	if !m.isFolded("e1") {
		t.Error("e1: explicit fold should be honoured")
	}
}

func TestEffectiveFolds(t *testing.T) {
	m := newFoldTestModel()
	m.foldOverrides["e2"] = false
	m.foldOverrides["m1"] = true

	folds := m.effectiveFolds(map[string]bool{"e2": true, "e3": true})

	want := map[string]bool{"e2": false, "e3": true, "m1": true}
	if len(folds) != len(want) {
		t.Fatalf("got %v, want %v", folds, want)
	}
	for id, expected := range want {
		if folds[id] != expected {
			t.Errorf("%s: got %v, want %v", id, folds[id], expected)
		}
	}
}

func TestSetFoldRecordsOverrideAndPendingSelection(t *testing.T) {
	m := newFoldTestModel()

	if cmd := m.setFold("e1", true); cmd == nil {
		t.Fatal("setFold should return a reload command")
	}

	if !m.foldOverrides["e1"] {
		t.Error("e1: expected fold override to be recorded")
	}
	if len(m.pendingSelectIDs) != 1 || m.pendingSelectIDs[0] != "e1" {
		t.Errorf("pendingSelectIDs = %v, want [e1]", m.pendingSelectIDs)
	}
}

func TestEffectiveFoldsSuspendedWhileFiltering(t *testing.T) {
	m := newFoldTestModel()
	m.foldOverrides["m1"] = true
	m.foldsSuspended = true

	// Folded children must be present in the list for the filter to match them
	if folds := m.effectiveFolds(map[string]bool{"e2": true}); len(folds) != 0 {
		t.Errorf("expected no folds while filtering, got %v", folds)
	}

	m.foldsSuspended = false
	if folds := m.effectiveFolds(map[string]bool{"e2": true}); !folds["e2"] || !folds["m1"] {
		t.Errorf("expected folds to return once filtering ends, got %v", folds)
	}
}

func TestRestoreSelectionKeepsOverrides(t *testing.T) {
	m := newFoldTestModel()
	m.pendingSelectIDs = []string{"t1", "e1", "m1"}

	// No items are loaded, so nothing matches; the request must still be cleared
	m.restoreSelection()

	if m.pendingSelectIDs != nil {
		t.Errorf("pendingSelectIDs = %v, want nil", m.pendingSelectIDs)
	}
}
