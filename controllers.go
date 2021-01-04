package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Signin todo doc
func Signin(w http.ResponseWriter, r *http.Request) {
	hm := NewHTTPMsg(w, *r)

	creds, err := hm.auth()

	if err != nil {
		log.Println(err)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := newClaims(creds, expirationTime)

	token, err := hm.tokenString(claims)
	if err != nil {
		log.Println(err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expirationTime,
	})
}

// Welcome todo doc
func Welcome(w http.ResponseWriter, r *http.Request) {
	hm := NewHTTPMsg(w, *r)

	claims, err := hm.validate()
	if err != nil {
		return
	}

	w.Write([]byte(fmt.Sprintf("Access granted to %s", claims.Username)))
}

// Refresh todo doc
func Refresh(w http.ResponseWriter, r *http.Request) {
	hm := NewHTTPMsg(w, *r)
	claims, err := hm.validate()

	if err != nil {
		log.Println(err)
		return
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
