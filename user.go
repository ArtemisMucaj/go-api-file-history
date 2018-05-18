package main

import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"

type User struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Mail      string
	Password  []byte
	Validated bool
	Token     Token
	Secret    []byte
}

func NewUserFromQuery(session *mgo.Session, query bson.M) (*User, error) {
	user := User{}
	collection := session.DB("database").C("user")
	err := collection.Find(query).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func NewUserFromToken(session *mgo.Session, t Token) (*User, error) {
	user := User{}
	collection := session.DB("database").C("user")
	err := collection.Find(bson.M{"token": t}).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) Save(session *mgo.Session) error {
	collection := session.DB("database").C("user")
	err := collection.Insert(bson.M{"mail": u.Mail,
		"password": u.Password})
	if err != nil {
		return err
	}
	return nil
}

func (u *User) SetToken(session *mgo.Session, t Token) error {
	collection := session.DB("database").C("user")
	err := collection.Update(bson.M{"mail": u.Mail},
		bson.M{"$set": bson.M{"token": t}})
	return err
}

func (u *User) UpdateFileReference(session *mgo.Session, uid string, commit string) error {
	collection := session.DB("database").C("user")
	err := collection.Update(bson.M{"mail": u.Mail},
		bson.M{"$set": bson.M{"files." + uid + ".commit": commit}})
	return err
}
