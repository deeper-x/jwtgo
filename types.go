package main

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// Credentials create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Claims todo doc
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// HTTPMsg todo doc
type HTTPMsg struct {
	Req  http.Request
	ResW http.ResponseWriter
}
