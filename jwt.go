package main

import "fmt"
import "time"
import "errors"
import "gopkg.in/dgrijalva/jwt-go.v3"

type Token string

func ValidateToken(t Token, secret []byte) bool {
	token, err := jwt.Parse(string(t), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return false
	}
	if token != nil && token.Valid {
		return true
	} else {
		return false
	}
}

func GenerateToken(mail string, secret []byte) (Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"mail": mail,
		"exp": time.Now().Add(time.Hour * 24).Unix()})
	tk, err := token.SignedString(secret)
	if err != nil {
		return Token(""), errors.New("Failed to sign token")
	}
	return Token(tk), nil
}
