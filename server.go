package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Username string
	Password string
}

type Err struct {
	Nope string
}

var dataLogin Login

func database(username string, email string, password string) {

	database, _ :=
		sql.Open("sqlite3", "data.db")
	// fmt.Println("openData base")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, username TEXT, email TEXT, password TEXT)")
	statement.Exec()
	statement, _ =
		database.Prepare("INSERT INTO people (username, email, password) VALUES (?, ?, ?)")
	statement.Exec(username, email, password)
	// fmt.Println("préparation de l'ajout")
	// fmt.Println(username)
	// fmt.Println(email)
	// fmt.Println(password)
	rows, _ :=
		database.Query("SELECT id, username, email, password FROM people")
	var id int
	fmt.Println("select data")
	for rows.Next() {
		rows.Scan(&id, &username, &email, &password)
		// fmt.Println(strconv.Itoa(id) + ": " + username + " " + email + " " + password)
	}
	database.Close()

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func verifRegisterData(regUser string, regEmail string) bool {
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
		// fmt.Println(strconv.Itoa(id) + ": " + username + " " + email + " " + password)
		if regEmail == email || regUser == username {
			fmt.Println("nope")
			database.Close()
			return false
		}
	}
	// fmt.Println("data : ", username, email)
	// fmt.Println("paramètre :", regUser, regEmail)
	database.Close()
	return true
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
		// fmt.Println(strconv.Itoa(id) + ": " + username + " " + email + " " + password)
		if logUser == username && CheckPasswordHash(logPassword, password) {
			fmt.Println("correct password")
			database.Close()
			return true
		}
	}
	database.Close()
	fmt.Println("wrong password")
	return false
}

func LoginHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")

	dataLogin.Username = r.FormValue("username")
	dataLogin.Password = r.FormValue("password")
	if dataLogin.Password != "" && dataLogin.Username != "" {
		if login(dataLogin.Username, dataLogin.Password) {
			http.Redirect(w, r, "/account", http.StatusSeeOther)
		}

	}
	tpl := template.Must(template.ParseFiles("assets/index.html"))

	err := tpl.Execute(w, dataLogin)
	if err != nil {
		log.Fatal(err)
	}
}

func RegisterHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Register")

	var regData Err
	userName := r.FormValue("username")
	email := r.FormValue("email")
	password, _ := HashPassword(r.FormValue("password"))

	// fmt.Println(userName)
	// fmt.Println(email)
	// fmt.Println(password)

	if userName != "" && email != "" && password != "" {
		if verifRegisterData(userName, email) {
			database(userName, email, password)

			fmt.Println("retourvers la page main")
			http.Redirect(w, r, "/main", http.StatusSeeOther)
		} else {
			regData.Nope = "your email or your username are already used ! chacal"
		}

	}

	tpl := template.Must(template.ParseFiles("assets/connect.html"))

	err := tpl.Execute(w, regData)
	if err != nil {
		log.Fatal(err)
	}
}

func account(w http.ResponseWriter, r *http.Request) {

	tpl := template.Must(template.ParseFiles("assets/signIn.html"))

	err := tpl.Execute(w, dataLogin)
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
