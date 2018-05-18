package main

import "strings"
import "context"
import "net/http"
import "gopkg.in/mgo.v2"

func (ctx *AppContext) SetDatabase(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := ctx.session.Copy()
		defer session.Close()
		this := context.WithValue(r.Context(), "session", session)
		inner.ServeHTTP(w, r.WithContext(this))
	})
}

func (ctx *AppContext) Authenticate(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token Token
		args, ok := r.Header["Authorization"]
		if ok && len(args) >= 1 {
			token = Token(strings.TrimPrefix(args[0], "Bearer "))
			session := r.Context().Value("session").(*mgo.Session)
			user, err := NewUserFromToken(session, token)
			if err == nil && ValidateToken(token, user.Secret) {
				this := context.WithValue(r.Context(), "user", user)
				inner.ServeHTTP(w, r.WithContext(this))
				return
			}
		}
		http.Error(w, http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)
		return
	})
}
