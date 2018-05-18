package main

import "errors"
import "net/http"

func GetMailAndPassword(req *http.Request) (string, string, error) {
	req.ParseForm()
	mail := req.Form.Get("mail")
	password := req.Form.Get("password")
	if mail != "" && password != "" {
		return mail, password, nil
	} else {
		return "", "", errors.New("Missing argument")
	}
}
