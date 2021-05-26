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

type Login struct {
	Username string
	Password string
}

var DataLogin Login

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
	database.Close()

}
func login(logUser string, logPassword string) bool {
	var username string
	var email string
	var password string
	var id int

	database, _ :=
		sql.Open("sqlite3", "data.db")
	rows, _ :=
		database.Query("SELECT id, username, email, password FROM people")
	for rows.Next() {
		rows.Scan(&id, &username, &email, &password)
		fmt.Println(strconv.Itoa(id) + ": " + username + " " + email + " " + password)
		if logUser == username && logPassword == password {
			fmt.Println("ok log id =", id)
			return true
		}
	}
	database.Close()
	return false
}

func LoginHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HOM")
	var data Login

	data.Username = r.FormValue("username")
	data.Password = r.FormValue("password")
	if data.Password != "" && data.Username != "" {
		if login(data.Username, data.Password) {
			http.Redirect(w, r, "/account", http.StatusSeeOther)
		}

	}
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
		database(userName, email, password)
		http.Redirect(w, r, "/main", http.StatusSeeOther)

	}

	data := "test"
	tpl := template.Must(template.ParseFiles("assets/connect.html"))

	err := tpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func account(w http.ResponseWriter, r *http.Request) {
	data := ""
	tpl := template.Must(template.ParseFiles("assets/signIn.html"))

	err := tpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.Handle("/", fs)
	http.HandleFunc("/main", LoginHandle)
	http.HandleFunc("/connect", RegisterHandle)
	http.HandleFunc("/account", account)
	http.ListenAndServe(":8080", nil)
}
