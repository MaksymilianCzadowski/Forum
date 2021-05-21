package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func mainHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HOM")
	data := "test"
	tpl := template.Must(template.ParseFiles("assets/index.html"))

	err := tpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.Handle("/", fs)
	http.HandleFunc("/main", mainHandle)
	http.ListenAndServe(":8080", nil)
}
