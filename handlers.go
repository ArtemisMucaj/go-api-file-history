package main

import "time"
import "net/http"
import "encoding/json"
import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"
import "golang.org/x/crypto/bcrypt"
import "github.com/satori/go.uuid"

func (ctx *AppContext) Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/views/index.html")
}

func (ctx *AppContext) Dev(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/views/dev.html")
}

func (ctx *AppContext) Login(w http.ResponseWriter, r *http.Request) {
	// Retrieve `mail` and `password` from http.Request pointer
	mail, password, err := GetMailAndPassword(r)
	if err != nil {
		http.Error(w, "Password incorrect", http.StatusUnauthorized)
		return
	}
	// Retrieve mgo.Session from the request's context then
	// load `user` using the `mail`
	session := r.Context().Value("session").(*mgo.Session)
	user, err := NewUserFromQuery(session, bson.M{"mail": mail})
	if err != nil {
		http.Error(w, "Password incorrect", http.StatusUnauthorized)
		return
	}
	// Verify that the password is correct (usually using a
	// slow hashing algorithm -> `bcrypt` is a good choice)
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		http.Error(w, "Password incorrect", http.StatusUnauthorized)
		return
	}
	if ValidateToken(user.Token, user.Secret) {
		w.Write([]byte(user.Token))
		return
	}
	// The `token` in database is no longer valid
	// We will generate and save a new one
	token, err := GenerateToken(user.Mail, user.Secret)
	if err != nil {
		http.Error(w, "Password incorrect", http.StatusUnauthorized)
		return
	}
	err = user.SetToken(session, token)
	if err != nil {
		http.Error(w, "Password incorrect", http.StatusUnauthorized)
		return
	}
	w.Write([]byte(token))
}

func (ctx *AppContext) Signin(w http.ResponseWriter, r *http.Request) {
	mail, password, err := GetMailAndPassword(r)
	if err != nil {
		http.Error(w, "Missing argument", http.StatusUnauthorized)
		return
	}
	session := r.Context().Value("session").(*mgo.Session)
	user := User{Mail: mail, Password: []byte(password)}
	err = user.Save(session)
	if err != nil {
		http.Error(w, "An unexpected error occurred", http.StatusUnauthorized)
		return
	}
}

func (ctx *AppContext) CreateFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uid := uuid.NewV4().String()
	commit := uuid.NewV4().String()
	filename := r.Form.Get("filename")
	branch := "master"
	message := r.Form.Get("message")
	data := r.Form.Get("data")
	if uid == "" || branch == "" || message == "" || data == "" ||
		filename == "" || commit == "" {
		http.Error(w, "Missing parameter", http.StatusUnauthorized)
		return
	}
	session := r.Context().Value("session").(*mgo.Session)
	user := r.Context().Value("user").(*User)
	timestamp := time.Now().Unix()
	file := File{Owners: []bson.ObjectId{user.ID}, UID: uid, Commit: commit, Filename: filename,
		Branch: branch, Message: message, Data: []byte(data), Timestamp: timestamp}
	err := file.Save(session)
	if err != nil {
		http.Error(w, "Unexpected error (could not save)", http.StatusUnauthorized)
		return
	}
	err = user.UpdateFileReference(session, file.UID, file.Commit)
	if err != nil {
		http.Error(w, "Unexpected error (could not update file reference)",
			http.StatusUnauthorized)
		return
	}
}

func (ctx *AppContext) SaveFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uid := r.Form.Get("uid")
	parent := r.Form.Get("parent")
	branch := r.Form.Get("branch")
	message := r.Form.Get("message")
	data := r.Form.Get("data")
	if uid == "" || parent == "" ||
		branch == "" || message == "" || data == "" {
		http.Error(w, "Missing parameter", http.StatusUnauthorized)
		return
	}
	session := r.Context().Value("session").(*mgo.Session)
	user := r.Context().Value("user").(*User)
	father, err := NewFileFromQuery(session, bson.M{"uid": uid, "parent": parent,
		"owners": user.ID})
	if err != nil {
		http.Error(w, "Unauthorized query", http.StatusUnauthorized)
		return
	}
	timestamp := time.Now().Unix()
	file := File{Owners: father.Owners, UID: uid, Filename: father.Filename,
		Branch: branch, Message: message, Parent: parent,
		Data: []byte(data), Timestamp: timestamp}
	err = file.Save(session)
	if err != nil {
		http.Error(w, "Unexpected error (could not save)", http.StatusUnauthorized)
		return
	}
	err = user.UpdateFileReference(session, file.UID, file.Commit)
	if err != nil {
		http.Error(w, "Unexpected error (could not update file reference)",
			http.StatusUnauthorized)
		return
	}
}

func (ctx *AppContext) GetFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	commit := r.Form.Get("commit")
	if commit == "" {
		http.Error(w, "Missing parameter", http.StatusUnauthorized)
		return
	}
	session := r.Context().Value("session").(*mgo.Session)
	user := r.Context().Value("user").(*User)
	file, err := NewFileFromQuery(session, bson.M{"commit": commit,
		"owners": user.ID})
	if err != nil {
		http.Error(w, "Unexpected error (could not get file)", http.StatusUnauthorized)
		return
	}
	json, err := json.Marshal(file)
	if err != nil {
		http.Error(w, "Unexpected error (could not serialize to json)", http.StatusUnauthorized)
		return
	}
	w.Write([]byte(json))
}

func (ctx *AppContext) GetFileHistory(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uid := r.Form.Get("uid")
	if uid == "" {
		http.Error(w, "Missing parameter", http.StatusUnauthorized)
		return
	}
	session := r.Context().Value("session").(*mgo.Session)
	user := r.Context().Value("user").(*User)
	files, err := NewFilesFromQuery(session, bson.M{"uid": uid,
		"owners": user.ID})
	if err != nil {
		http.Error(w, "Unexpected error (could not load files)", http.StatusUnauthorized)
		return
	}
	result, err := MakeHistoryTree(files)
	if err != nil {
		http.Error(w, "Unexpected error (could not build file history)", http.StatusUnauthorized)
		return
	}
	json, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Unexpected error (could not serialize to json)", http.StatusUnauthorized)
		return
	}
	w.Write([]byte(json))
}
