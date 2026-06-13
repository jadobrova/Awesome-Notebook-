package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed index.html
var content embed.FS

func main() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/load", loadHandler)
	log.Println("Go блокнот запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFS(content, "index.html")
	tmpl.Execute(w, nil)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	filename := r.FormValue("filename")
	content := r.FormValue("content")
	if filename == "" {
		w.Write([]byte("error: no filename"))
		return
	}
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		w.Write([]byte("error: " + err.Error()))
		return
	}
	w.Write([]byte("ok"))
}

func loadHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Error(w, "missing file", http.StatusBadRequest)
		return
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		w.Write([]byte(""))
		return
	}
	w.Write(data)
}
