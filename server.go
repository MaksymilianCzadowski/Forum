package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Username string
	Password string
}

type Err struct {
	Nope         string
	WrongPass    string
	UserNotFound string
}

type Post struct {
	Comment    string
	TagCod     string
	TagMusic   string
	TagArt     string
	TagSport   string
	TagFashion string
	TagFood    string
	TagCinema  string
	TagCars    string
}

type PrintPost struct {
	Username string
	Id       int
	Comment  string
}

type postMessage struct {
	id      int
	image   string
	pseudo  string
	titre   string
	message string
}

type oui struct {
	Comment []postMessage
}

var allData []PrintPost
var errToSend Err
var store = sessions.NewCookieStore([]byte("mysession"))
var dataLogin Login
var tag Post
var database, _ = sql.Open("sqlite3", "data.db")

func createTablePeople() {
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, username TEXT, email TEXT, password TEXT)")
	defer statement.Close()
	statement.Exec()
}

func createTablePost() {
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS post (id INTEGER PRIMARY KEY, username TEXT, comment TEXT)")
	defer statement.Close()
	statement.Exec()
}

func addUser(username string, email string, password string) {

	statement, _ :=
		database.Prepare("INSERT INTO people (username, email, password) VALUES (?, ?, ?)")
	defer statement.Close()
	statement.Exec(username, email, password)

	lastRow := database.QueryRow("SELECT * FROM user WHERE id=(SELECT max(id) FROM user)")
	var id int
	lastRow.Scan(&id, &email, &username, &password)
	fmt.Println("NEW USER : ", username)

}

func addPost(username string) {
	fmt.Println(username)
	fmt.Println(tag.Comment)
	// fmt.Println(tag.TagCars)

	statement, _ :=
		database.Prepare("INSERT INTO post (username, comment) VALUES (?, ?)")
	// defer statement.Close()
	statement.Exec(username, tag.Comment)

	rows, _ :=
		database.Query("SELECT id, username, comment FROM post")
	var id int
	fmt.Println("select data")
	for rows.Next() {
		rows.Scan(&id, &username, &tag.Comment)
		fmt.Println(strconv.Itoa(id) + ": " + username + " " + tag.Comment)
	}

}

// func databaseComment(id int, username string, newComment string) {

// 	database, _ :=
// 		sql.Open("sqlite3", "data.db")
// 	statement, _ :=
// 		database.Prepare("CREATE TABLE IF NOT EXISTS comment (id INTEGER , username TEXT, newComment TEXT)")
// 	statement.Exec()
// 	statement, _ =
// 		database.Prepare("INSERT INTO comment (id, username, newComment) VALUES (?, ?, ?)")
// 	fmt.Println("ici")
// 	statement.Exec(id, username, newComment)
// 	rows, _ :=
// 		database.Query("SELECT id, username, newComment FROM comment")
// 	var test []string
// 	for rows.Next() {
// 		rows.Scan(&id, &username, &newComment)
// 		test = append(test, strconv.Itoa(id)+": "+username+" "+newComment+"\n")
// 	}

// }

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
	var result = true

	rows, _ :=
		database.Query("SELECT id, username, email, password FROM people")
	for rows.Next() {
		rows.Scan(&id, &username, &email, &password)
		// fmt.Println(strconv.Itoa(id) + ": " + username + " " + email + " " + password)
		if regEmail == email || regUser == username {
			fmt.Println("nope")
			result = false
		}
	}
	// fmt.Println("data : ", username, email)
	// fmt.Println("paramètre :", regUser, regEmail)
	return result
}

func login(logUser string, logPassword string) int {
	var username string
	var email string
	var password string
	var id int
	var result = 0

	rows, _ :=
		database.Query("SELECT id, username, email, password FROM people")
	for rows.Next() {
		rows.Scan(&id, &username, &email, &password)
		// fmt.Println(strconv.Itoa(id) + ": " + username + " " + email + " " + password)
		if logUser == username && CheckPasswordHash(logPassword, password) {
			result = 1
		} else if logUser == username && !(CheckPasswordHash(logPassword, password)) {
			result = 2
		}
	}
	return result
}

func logoutHandle(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mysession")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/main", http.StatusSeeOther)
}

func LoginHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")
	dataLogin.Username = r.FormValue("username")
	dataLogin.Password = r.FormValue("password")
	errToSend.WrongPass = ""
	errToSend.UserNotFound = ""
	tpl := template.Must(template.ParseFiles("assets/index.html"))

	if dataLogin.Password != "" && dataLogin.Username != "" {
		if login(dataLogin.Username, dataLogin.Password) == 1 {
			fmt.Println("correct password")
			session, _ := store.Get(r, "mysession")
			session.Values["username"] = dataLogin.Username
			session.Save(r, w)
			http.Redirect(w, r, "/account", http.StatusSeeOther)
		} else if login(dataLogin.Username, dataLogin.Password) == 2 {
			fmt.Println("Wrong password")
			errToSend.WrongPass = "Wrong password"
			tpl.Execute(w, errToSend)
			return
		} else {
			fmt.Println("user doesn't exist")
			errToSend.UserNotFound = "user doesn't exist"
			tpl.Execute(w, errToSend)
			return
		}

	}

	err := tpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllData() {
	var temp PrintPost

	rows, _ :=
		database.Query("SELECT id, username, comment FROM post")

	for rows.Next() {
		rows.Scan(&temp.Id, &temp.Username, &temp.Comment)
		allData = append(allData, temp)
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
			addUser(userName, email, password)
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
	getAllData()

	session, _ := store.Get(r, "mysession")
	userName := fmt.Sprintf("%v", session.Values["username"])
	data := map[string]interface{}{
		"username": userName,
		"post":     allData,
	}
	// fmt.Println(allData[0])

	tpl, _ := template.ParseFiles("assets/signIn.html")
	tpl.Execute(w, data)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

// func CreatePostHandle(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("POST")

// 	session, _ := store.Get(r, "mysession")
// 	username := fmt.Sprintf("%v", session.Values["username"])

// 	tag.Comment = r.FormValue("comment")
// 	tag.TagCod = r.FormValue("tagCod")
// 	tag.TagMusic = r.FormValue("tagMusic")
// 	tag.TagArt = r.FormValue("tagArt")
// 	tag.TagSport = r.FormValue("tagSport")
// 	tag.TagFashion = r.FormValue("tagFashion")
// 	tag.TagFood = r.FormValue("tagFood")
// 	tag.TagCinema = r.FormValue("tagCinema")
// 	tag.TagCars = r.FormValue("tagCars")

// 	if tag.Comment != "" {
// 		fmt.Printf("psot en cour\n")
// 		addPost(username)
// 	}

// 	tpl := template.Must(template.ParseFiles("assets/post.html"))
// 	tpl.Execute(w, nil)
// }

// func CommentHandle(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Connect")

// 	id := r.FormValue("id")
// 	userName := PrintPost.Username
// 	newComment := r.FormValue("newComment")
// 	intId, _ := strconv.Atoi(id)
// 	if newComment != "" {
// 		databaseComment(intId, userName, newComment)
// 		fmt.Println("tu sors de  wsh")
// 	}

// 	tpl := template.Must(template.ParseFiles("assets/comment.html"))

// 	data := oui{
// 		Comment: getCommentInfo(),
// 	}

// 	err := tpl.Execute(w, data)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// }

// func getCommentInfo() []postMessage {

// 	database, _ :=
// 		sql.Open("sqlite3", "data.db")
// 	rows, _ :=
// 		database.Query("SELECT id, username, newComment FROM comment")

// 	var _id int
// 	var test []postMessage
// 	var _pseudo string
// 	var _message string
// 	for rows.Next() {
// 		rows.Scan(&_id, &_pseudo, &_message)
// 		data := postMessage{
// 			id:      _id,
// 			pseudo:  _pseudo,
// 			message: _message,
// 		}
// 		test = append(test, data)
// 		fmt.Println(test)
// 	}
// 	return test
// }

func main() {
	createTablePeople()
	createTablePost()
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.Handle("/", fs)
	http.HandleFunc("/main", LoginHandle)
	http.HandleFunc("/logout", logoutHandle)
	// http.HandleFunc("/connect", RegisterHandle)
	http.HandleFunc("/account", PostHandle)
	http.HandleFunc("/CreateNewPost", CreatePostHandle)
	http.HandleFunc("/comment", CommentHandle)

	http.ListenAndServe(":8080", nil)
}
