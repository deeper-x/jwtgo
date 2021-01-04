package main

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

// NewHTTPMsg todo doc
func NewHTTPMsg(w http.ResponseWriter, r http.Request) HTTPMsg {
	return HTTPMsg{
		Req:  r,
		ResW: w,
	}
}

func (h HTTPMsg) validate() (*Claims, error) {
	claims := &Claims{}
	c, err := h.Req.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			h.ResW.WriteHeader(http.StatusUnauthorized)
			return claims, err
		}
		h.ResW.WriteHeader(http.StatusBadRequest)
		return claims, err
	}
	tknStr := c.Value

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			h.ResW.WriteHeader(http.StatusUnauthorized)
			return claims, err
		}
		h.ResW.WriteHeader(http.StatusBadRequest)
		return claims, err
	}

	if !tkn.Valid {
		h.ResW.WriteHeader(http.StatusUnauthorized)
		return claims, err
	}

	return claims, nil
}

func (h HTTPMsg) auth() (Credentials, error) {
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(h.Req.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		h.ResW.WriteHeader(http.StatusBadRequest)
		return creds, err
	}

	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		h.ResW.WriteHeader(http.StatusUnauthorized)
		return creds, err
	}

	return creds, nil
}

func newClaims(creds Credentials, expirationTime time.Time) *Claims {
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	return claims
}

func (h HTTPMsg) tokenString(c *Claims) (string, error) {
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		h.ResW.WriteHeader(http.StatusInternalServerError)
		return tokenString, err
	}

	return tokenString, nil
}
