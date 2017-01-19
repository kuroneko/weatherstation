package main

import (
	"net/http"
)

func init() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", http.RedirectHandler("/static/", http.StatusFound))
}
