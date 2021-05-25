package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func database(username string, email string, password string) {

	database, _ :=
		sql.Open("sqlite3", "data.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, username TEXT, email TEXT, password TEXT)")
	statement.Exec()
	statement, _ =
		database.Prepare("INSERT INTO people (username, email, password) VALUES (?, ?, ?)")
	statement.Exec(username, email, password)
	rows, _ :=
		database.Query("SELECT id, username, email, password FROM people")
	var id int
	for rows.Next() {
		rows.Scan(&id, &username, &email, &password)
		fmt.Println(strconv.Itoa(id) + ": " + username + " " + email + " " + password)
	}

}

func mainHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HOM")
	// userName := r.FormValue("name")
	// password := r.FormValue("password")
	data := "test"
	tpl := template.Must(template.ParseFiles("assets/index.html"))

	err := tpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func RegisterHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connect")

	userName := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if userName != "" && email != "" && password != "" {
		fmt.Println("tu rentre de dans wsh")
		database(userName, email, password)
	}

	data := "test"
	tpl := template.Must(template.ParseFiles("assets/connect.html"))

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

	http.HandleFunc("/connect", RegisterHandle)
	http.ListenAndServe(":8080", nil)
}
