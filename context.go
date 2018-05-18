package main

import "time"
import "gopkg.in/mgo.v2"

type AppContext struct {
	session *mgo.Session
}

func NewAppContext() (*AppContext, error) {
	session, err := mgo.DialWithTimeout("mongodb:27017/",
		time.Duration(5*time.Second))
	if err != nil {
		return nil, err
	}
	return &AppContext{session: session}, nil
}

func (ctx *AppContext) Close() {
	ctx.session.Close()
}
