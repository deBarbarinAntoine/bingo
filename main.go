package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"time"
	
	"BinGo/binder"
)

type Publication struct {
	Title string    `multipart:"title"`
	Score float64   `multipart:"score"`
	Date  time.Time `multipart:"date"`
}

type Info struct {
	Id           uint          `multipart:"id"`
	Publications []Publication `multipart:"publications"`
}

type Data struct {
	Name    string                `multipart:"name"`
	Age     int                   `multipart:"age"`
	IsAdmin bool                  `multipart:"is_admin"`
	Score   float64               `multipart:"score"`
	Info    Info                  `multipart:"info"`
	File    *multipart.FileHeader `multipart:"file"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	var data = &Data{}
	query, err := binder.NewMultipartForm(data, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = query.Fetch()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if data.File != nil {
		fmt.Printf("%d bytes of file uploaded: %s\n", data.File.Size, data.File.Filename)
	} else {
		fmt.Println("No file uploaded")
	}
	
	w.Header().Set("Content-Type", "application/json")
	
	json.NewEncoder(w).Encode(data)
}

//go:generate go run github.com/debarbarinantoine/go-enum-generate
func main() {
	http.HandleFunc("POST /", handler)
	log.Fatal(http.ListenAndServe(":8008", nil))
}
