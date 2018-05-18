package main

import "testing"

func TestOrderedMakeHistoryTree(t *testing.T) {
	files := []File{}
	// First file
	f1 := File{UID: "unittest",
		Commit: "first_commit"}
	files = append(files, f1)
	// Second file
	f2 := File{UID: "unitteft",
		Commit: "second_commit"}
	files = append(files, f2)
	// Third file
	f3 := File{UID: "unittest",
		Parent: "first_commit",
		Commit: "third_commit"}
	files = append(files, f3)
	history, _ := MakeHistoryTree(files)
	if history.Sons[0].Commit != "first_commit" {
		t.Error("Not the right commit (first)")
	}
	if history.Sons[1].Commit != "second_commit" {
		t.Error("Not the right commit (second)")
	}
	if history.Sons[0].Sons[0].Commit != "third_commit" {
		t.Error("Not the right commit (third)")
	}
}

func TestUnorderedMakeHistoryTree(t *testing.T) {
	files := []File{}
	// First file
	f1 := File{UID: "unittest",
		Parent: "first_commit",
		Commit: "third_commit"}
	files = append(files, f1)
	// Second file
	f2 := File{UID: "unittest",
		Commit: "first_commit"}
	files = append(files, f2)
	// Third file
	f3 := File{UID: "unittest",
		Commit: "second_commit"}
	files = append(files, f3)
	history, _ := MakeHistoryTree(files)
	if history.Sons[0].Commit != "first_commit" {
		t.Error("Not the right commit (first)")
	}
	if history.Sons[1].Commit != "second_commit" {
		t.Error("Not the right commit (second)")
	}
	if history.Sons[0].Sons[0].Commit != "third_commit" {
		t.Error("Not the right commit (third)")
	}
}

func TestComplexUnorderedMakeHistoryTree(t *testing.T) {
	files := []File{}
	// First file
	f1 := File{UID: "unittest",
		Parent: "first_commit",
		Commit: "third_commit"}
	files = append(files, f1)
	// Second file
	f2 := File{UID: "unittest",
		Commit: "first_commit"}
	files = append(files, f2)
	// Third file
	f3 := File{UID: "unittest",
		Commit: "second_commit"}
	files = append(files, f3)
	// Fourth file
	f4 := File{UID: "unittest",
		Parent: "first_commit",
		Commit: "fourth_commit"}
	files = append(files, f4)
	history, _ := MakeHistoryTree(files)
	if history.Sons[0].Commit != "first_commit" {
		t.Error("Not the right commit (first)")
	}
	if history.Sons[1].Commit != "second_commit" {
		t.Error("Not the right commit (second)")
	}
	if history.Sons[0].Sons[0].Commit != "third_commit" {
		t.Error("Not the right commit (third)")
	}
	if history.Sons[0].Sons[1].Commit != "fourth_commit" {
		t.Error("Not the right commit (fourth)")
	}
}
