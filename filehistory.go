package main

import "errors"
import "gopkg.in/mgo.v2/bson"

type Children []*History

type History struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	UID    string        `bson:"uid"`
	Parent string
	Commit string
	Sons   Children
}

func (f *History) Find(commit string) (*History, error) {
	if f.Commit == commit {
		return f, nil
	}
	version, err := f.Sons.Find(commit)
	return version, err
}

func (f *History) IsChild(h *History) bool {
	for _, child := range f.Sons {
		if h.Commit == child.Parent {
			return true
		}
	}
	return false
}

func (f *History) ChangeParent(h *History) error {
	children := Children{}
	for _, child := range f.Sons {
		if h.Commit == child.Parent {
			h.Sons = append(h.Sons, child)
		} else {
			children = append(children, child)
		}
	}
	f.Sons = children
	f.Sons = append(f.Sons, h)
	return nil
}

func (f *History) Add(file File) error {
	h, _ := MakeHistory(file)
	if f.IsChild(h) {
		f.ChangeParent(h)
		return nil
	}
	f.Sons = append(f.Sons, h)
	return nil
}

func (f Children) Find(commit string) (*History, error) {
	for _, elt := range f {
		version, error := elt.Find(commit)
		if error != nil {
			continue
		}
		return version, nil
	}
	return nil, errors.New("Not found")
}

func (f Children) Add(file File) error {
	h, _ := MakeHistory(file)
	f = append(f, h)
	return nil
}

func MakeHistory(file File) (*History, error) {
	h := History{}
	h.ID = file.ID
	h.UID = file.UID
	h.Parent = file.Parent
	h.Commit = file.Commit
	h.Sons = Children{}
	return &h, nil
}

func MakeHistoryTree(files []File) (*History, error) {
	h := History{}
	for _, file := range files {
		father, err := h.Find(file.Parent)
		if err != nil {
			h.Add(file)
			continue
		}
		father.Add(file)
	}
	return &h, nil
}
