package main

import "time"
import "testing"
import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"
import "github.com/satori/go.uuid"

func TestNewFileFromQuery(t *testing.T) {
	session, _ := mgo.DialWithTimeout("mongodb:27017/",
		time.Duration(5 * time.Second))
	// Insert fake elements
	timestamp := time.Now().Unix()
	collection := session.DB("database").C("file")
	collection.Insert(bson.M{"uid": "integrationtest",
		"filename": "fakefile", "commit": "first_commit",
		"timestamp": timestamp})
	// Test
	file, _ := NewFileFromQuery(session, bson.M{"uid": "integrationtest"})
	if file == nil {
		t.Error("File is nil")
	}
	if file != nil && string(file.UID) != "integrationtest" {
		t.Error("Hash value is wrong")
	}
	if file != nil && file.Filename != "fakefile" {
		t.Error("Filename value is wrong")
	}
	if file != nil && file.Commit != "first_commit" {
		t.Error("Commit value is wrong")
	}
	// Clean up database
	collection.RemoveAll(bson.M{"uid" : "integrationtest"})
}

func TestNewFileFromQueryByOwner(t *testing.T) {
	session, _ := mgo.DialWithTimeout("mongodb:27017/",
		time.Duration(5 * time.Second))
	// Insert fake elements
	timestamp := time.Now().Unix()
	collection := session.DB("database").C("file")
	collection.Insert(bson.M{"uid": "integrationtest",
		"filename": "fakefile", "commit": "first_commit",
		"timestamp": timestamp, "owners" : []interface{}{"firstuser", "seconduser"}})
	// Test
	file, _ := NewFileFromQuery(session, bson.M{"uid": "integrationtest",
		"owners": "firstuser"})
	if file == nil {
		t.Error("File is nil")
	}
	if file != nil && string(file.UID) != "integrationtest" {
		t.Error("Hash value is wrong")
	}
	if file != nil && file.Filename != "fakefile" {
		t.Error("Filename value is wrong")
	}
	if file != nil && file.Commit != "first_commit" {
		t.Error("Commit value is wrong")
	}
	// Clean up database
	collection.RemoveAll(bson.M{"filename" : "fakefile"})
}

func TestNewFilesFromQuery(t *testing.T) {
	session, _ := mgo.DialWithTimeout("mongodb:27017/",
		time.Duration(5 * time.Second))
	// Insert fake elements
	timestamp := time.Now().Unix()
	collection := session.DB("database").C("file")
	collection.Insert(bson.M{"uid": "integrationtest",
		"filename": "fakefile", "commit": "first_commit",
		"timestamp": timestamp})
	collection.Insert(bson.M{"uid": "integrationtest",
		"filename": "fakefile", "parent": "first_commit",
		"commit": "second_commit", "timestamp": timestamp + 1})
	// Test
	files, _ := NewFilesFromQuery(session, bson.M{"uid": "integrationtest"})
	if len(files) != 2 {
		t.Error("Length of files array is not equal 2 but", len(files))
	}
	if string(files[0].UID) != "integrationtest" {
		t.Error("Hash value is wrong")
	}
	if files[0].Filename != "fakefile" {
		t.Error("Filename value is wrong")
	}
	if files[0].Commit != "first_commit" {
		t.Error("Commit value is wrong")
	}
	if string(files[1].UID) != "integrationtest" {
		t.Error("Hash value is wrong")
	}
	if files[1].Filename != "fakefile" {
		t.Error("Filename value is wrong")
	}
	if files[1].Commit != "second_commit" {
		t.Error("Commit value is wrong")
	}
	if files[1].Parent != "first_commit" {
		t.Error("Commit value is wrong")
	}
	// Clean up database
	collection.RemoveAll(bson.M{"filename" : "fakefile"})
}

func TestCreate(t *testing.T) {
	session, _ := mgo.DialWithTimeout("mongodb:27017/",
		time.Duration(5 * time.Second))
	collection := session.DB("database").C("file")
	timestamp := time.Now().Unix()
	file := File{Filename: "spur-gear", Parent: "test",
		Branch: "master", Message: "First commit",
		Data: []byte("empty"), Timestamp: timestamp}
	file.Create(session)
	// Test
	res, _ := NewFileFromQuery(session, bson.M{"uid": file.UID})
	if res.Parent != "" {
		t.Error("Wrong parent field")
	}
	if res.Branch != "master" {
		t.Error("Wrong branch field")
	}
	if res.Message != "First commit" {
		t.Error("Wrong commit message")
	}
	if res.Timestamp != timestamp {
		t.Error("Wrong timestamp value")
	}
	_, err := uuid.FromString(res.Commit)
    if err != nil {
        t.Error("Wrong uuid value in `Commit`")
    }
	// Clean up database
	collection.RemoveAll(bson.M{"filename" : "spur-gear"})
}

func TestSave(t *testing.T) {
	session, _ := mgo.DialWithTimeout("mongodb:27017/",
		time.Duration(5 * time.Second))
	collection := session.DB("database").C("file")
	timestamp := time.Now().Unix()
	file := File{UID: "integrationtest",
		Filename: "spur-gear", Parent: "test",
		Branch: "master", Message: "First commit",
		Data: []byte("empty"), Timestamp: timestamp}
	file.Save(session)
	// Test
	res, _ := NewFileFromQuery(session, bson.M{"uid": "integrationtest"})
	if res.Parent != "test" {
		t.Error("Wrong parent field")
	}
	if res.Branch != "master" {
		t.Error("Wrong branch field")
	}
	if res.Message != "First commit" {
		t.Error("Wrong commit message")
	}
	if res.Timestamp != timestamp {
		t.Error("Wrong timestamp value")
	}
	_, err := uuid.FromString(res.Commit)
    if err != nil {
        t.Error("Wrong uuid value in `Commit`")
    }
	// Clean up database
	collection.RemoveAll(bson.M{"filename" : "spur-gear"})
}