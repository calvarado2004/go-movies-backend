package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
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

// generateTokenPair generates a new access and refresh token pair.
func (j *Auth) generateTokenPair(user *jwtUser) (tokenPairs, error) {

	// Create a token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = j.Audience
	claims["iss"] = j.Issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"
	claims["exp"] = time.Now().UTC().Add(j.TokenExpiry).Unix()

	// Sign the token
	signedAccessToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return tokenPairs{}, err
	}

	// Create a refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	// Set the claims
	claimsRefresh := refreshToken.Claims.(jwt.MapClaims)
	claimsRefresh["sub"] = fmt.Sprint(user.ID)
	claimsRefresh["iat"] = time.Now().UTC().Unix()
	claimsRefresh["exp"] = time.Now().UTC().Add(j.TokenExpiry).Unix()

	// Sign the refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		return tokenPairs{}, err
	}

	// Create a tokenPairs struct and populate it with the tokens
	tokenPairs := tokenPairs{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	// Return the tokenPairs struct and nil as the error value
	return tokenPairs, nil

}

// getRefreshCookie generates a new refresh cookie.
func (j *Auth) getRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Value:    refreshToken,
		Path:     j.CookiePath,
		Domain:   j.CookieDomain,
		Expires:  time.Now().UTC().Add(j.RefreshExpiry),
		MaxAge:   int(j.RefreshExpiry.Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

// getExpiredRefreshCookie generates a new refresh cookie with an expired time.
func (j *Auth) getExpiredRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Value:    "",
		Path:     j.CookiePath,
		Domain:   j.CookieDomain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}
