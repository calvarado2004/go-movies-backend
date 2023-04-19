package main

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// Auth is a struct that holds the authentication configuration.
type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	CookieDomain  string
	CookieName    string
	CookiePath    string
}

// jwtUser returns a new Auth struct.
type jwtUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// tokenPairs is a struct that holds the access and refresh tokens.
type tokenPairs struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// tokenClaims is a struct that holds the claims for the access token.
type tokenClaims struct {
	jwt.RegisteredClaims
}
