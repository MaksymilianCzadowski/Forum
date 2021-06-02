package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

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

type postInfo struct {
	idx     int
	pseudo  string
	message string
}
type postMessage struct {
	id      int
	image   string
	pseudo  string
	titre   string
	message string
}

type oui struct {
	Poste []postMessage
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

func databasePost(username string, newPost string) {

	database, _ :=
		sql.Open("sqlite3", "data.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS post (id INTEGER PRIMARY KEY, username TEXT, newPost TEXT)")
	statement.Exec()
	statement, _ =
		database.Prepare("INSERT INTO post (username, newPost) VALUES (?, ?)")
	fmt.Println("ici")
	statement.Exec(username, newPost)
	fmt.Println("G TROUVER VOUS ETES NUL")
	rows, _ :=
		database.Query("SELECT id, username, newPost FROM post")
	var id int
	var test []string
	for rows.Next() {
		rows.Scan(&id, &username, &newPost)
		test = append(test, strconv.Itoa(id)+": "+username+" "+newPost+"\n")
	}

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

func PostHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connect")

	userName := r.FormValue("username")
	newPost := r.FormValue("newPost")

	if userName != "" && newPost != "" {
		databasePost(userName, newPost)
		fmt.Println("tu sors de  wsh")
		// http.Redirect(w, r, "/main", http.StatusSeeOther)
	}

	// var infoPost []string
	// infoPost = append(infoPost, userName)
	// infoPost = append(infoPost, newPost)
	// println(&infoPost)

	tpl := template.Must(template.ParseFiles("assets/signIn.html"))

	data := oui{
		Poste: getPostInfo(),
	}

	err := tpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}

}

func getPostInfo() []postMessage {

	database, _ :=
		sql.Open("sqlite3", "data.db")
	rows, _ :=
		database.Query("SELECT id, username, newPost FROM post")

	var _id int
	var test []postMessage
	var _pseudo string
	var _message string
	for rows.Next() {
		rows.Scan(&_id, &_pseudo, &_message)
		data := postMessage{
			id:      _id,
			pseudo:  _pseudo,
			message: _message,
		}
		test = append(test, data)
	}

	// fmt.Println(strconv.Itoa(id) + ": " + username + " " + newPost)
	// fmt.Println(test)
	return test

}

// func account(w http.ResponseWriter, r *http.Request) {

// 	tpl := template.Must(template.ParseFiles("assets/signIn.html"))

// 	err := tpl.Execute(w, dataLogin)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.Handle("/", fs)
	http.HandleFunc("/main", LoginHandle)
	http.HandleFunc("/connect", RegisterHandle)
	http.HandleFunc("/account", PostHandle)
	http.ListenAndServe(":8080", nil)
}
