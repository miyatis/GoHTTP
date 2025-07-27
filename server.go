package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
)

func handleError(w http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

func main() {
	var httpServer http.Server
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/submit", submitHandler)

	httpServer.Addr = ":18888"

	fmt.Println("サーバーが http://localhost:18888 で起動しました")
	log.Println(httpServer.ListenAndServe())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if handleError(w, err) {
		return
	}

	tmpl, err := template.ParseFiles("index.html")
	if handleError(w, err) {
		return
	}

	tmpl.Execute(w, nil)
	fmt.Println(string(dump))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if handleError(w, err) {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドのみ許可されています", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	author := r.FormValue("author")

	tmpl, err := template.ParseFiles("success.html")
	if handleError(w, err) {
		return
	}

	data := struct {
		Title  string
		Author string
	}{
		Title:  title,
		Author: author,
	}

	fmt.Println(string(dump))
	tmpl.Execute(w, data)
}

