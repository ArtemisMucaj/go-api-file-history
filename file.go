package main

import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"
import "github.com/satori/go.uuid"

type File struct {
	ID        bson.ObjectId   `json:"_id" bson:"_id,omitempty"`
	Owners    []bson.ObjectId `json:"owners" bson:"owners"`
	UID       string          `json:"uid" bson:"uid"`
	Filename  string          `json:"filename" bson:"filename"`
	Branch    string          `json:"branch" bson:"branch"`
	Message   string          `json:"message" bson:"message"`
	Parent    string          `json:"parent" bson:"parent"`
	Commit    string          `json:"commit" bson:"commit"`
	Data      []byte          `json:"data" bson:"data"`
	Timestamp int64           `json:"timestamp" bson:"timestamp"`
}

func NewFileFromQuery(session *mgo.Session, query bson.M) (*File, error) {
	file := File{}
	collection := session.DB("database").C("file")
	err := collection.Find(query).One(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func NewFilesFromQuery(session *mgo.Session, query bson.M) ([]File, error) {
	files := []File{}
	collection := session.DB("database").C("file")
	err := collection.Find(query).All(&files)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (f *File) Create(session *mgo.Session) error {
	collection := session.DB("database").C("file")
	f.UID = uuid.NewV4().String()
	f.Commit = uuid.NewV4().String()
	err := collection.Insert(bson.M{"owners": f.Owners,
		"uid": f.UID, "filename": f.Filename, "branch": f.Branch,
		"message": f.Message, "commit": f.Commit,
		"data": f.Data, "timestamp": f.Timestamp})
	if err != nil {
		return err
	}
	return nil
}

func (f *File) Save(session *mgo.Session) error {
	collection := session.DB("database").C("file")
	f.Commit = uuid.NewV4().String()
	err := collection.Insert(bson.M{"owners": f.Owners,
		"uid": f.UID, "filename": f.Filename, "branch": f.Branch,
		"message": f.Message, "parent": f.Parent,
		"commit": f.Commit, "data": f.Data,
		"timestamp": f.Timestamp})
	if err != nil {
		return err
	}
	return nil
}
