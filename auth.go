package main

import (
	"crypto/subtle"
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"log"
	"net/http"
	"time"
	pass "github.com/sethvargo/go-password/password"
)

const TOKEN_NAME = "MYSQL_CONSOLE_TOKEN"
var auth_secret = ""

func init() {
	auth_secret = pass.MustGenerate(64, 10, 10, false, false)
}

type Claims struct {
	jwt.StandardClaims
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	// t, err := template.ParseFS(templateFiles, "templates/login.html")
	t, err := template.ParseFiles("templates/base.html", "templates/login.html")
	if err != nil {
		log.Println("Error parsing login.html: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t.Execute(w, Env{
		Title: title,
	})
}

// TODO: handle redirects in js
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		u := r.FormValue("username")
		p := r.FormValue("password")

		if subtle.ConstantTimeCompare([]byte(u), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(p), []byte(password)) != 1 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		expiration := time.Now().Add(24 * time.Hour * 365)
		claims := Claims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expiration.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
		signed, err := token.SignedString([]byte(auth_secret))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:       TOKEN_NAME,
			Value:      signed,
			Expires:    expiration,
		})
	} else {
		loginPage(w, r)
	}
}

func verify(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(TOKEN_NAME)
		if err != nil {
			loginPage(w, r)
			return
		}

		var claims Claims
		_, err = jwt.ParseWithClaims(cookie.Value, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(auth_secret), nil
		})
		if err != nil {
			loginPage(w, r)
			return
		}

		next(w, r)
	})
}