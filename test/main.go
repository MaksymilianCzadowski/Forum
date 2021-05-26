// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"html/template"
// 	"log"
// 	"net/http"
// 	"strconv"
// )

// type infoData struct {
// 	id       int
// 	username string
// 	newPost    string
// 	password string
// 	newPost  string
// }

// func mainHandle(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("HOM")
// 	data := "test"
// 	tpl := template.Must(template.ParseFiles("assets/index.html"))
// 	err := tpl.Execute(w, data)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func main() {
// 	fmt.Println("TTTrrrrooooppppp ccccooooollllll ssaaaaaa mmmmmarccchhhhheeee :)))))))")
// 	fs := http.FileServer(http.Dir("assets"))
// 	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
// 	http.Handle("/", fs)
// 	http.HandleFunc("/main", mainHandle)
// 	http.HandleFunc("/post", PostHandle)
// 	http.ListenAndServe(":6060", nil)
// }

// func database(username string, newPost string) {

// 	database, _ :=
// 		sql.Open("sqlite3", "data.db")
// 	statement, _ :=
// 		database.Prepare("CREATE TABLE IF NOT EXISTS post (id INTEGER PRIMARY KEY, username TEXT, newPost TEXT)")
// 	statement.Exec()
// 	statement, _ =
// 		database.Prepare("INSERT INTO post (username, newPost) VALUES (?, ?)")
// 	statement.Exec(username, newPost)
// 	rows, _ :=
// 		database.Query("SELECT id, username, newPost FROM post")
// 	var id int
// 	for rows.Next() {
// 		rows.Scan(&id, &username, &newPost)
// 		fmt.Println(strconv.Itoa(id) + ": " + username + " " + newPost)
// 	}

// }

// func PostHandle(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("post")

// 	// userName := r.FormValue("username")
// 	// userName := "Chocolat"
// 	// newPost := r.FormValue("newPost")
// 	var postTab []string

// 	idata := infoData{
// 		username: "doof",
// 		newPost:  "Salut tout le monde.",
// 	}
// 	tpl := template.Must(template.ParseFiles("assets/test.html"))

// 	// if userName != "" && newPost != "" {
// 	fmt.Println("tu rentres dedans wsh")
// 	database(idata.username, idata.newPost)
// 	http.Redirect(w, r, "/main", http.StatusSeeOther)
// 	// }

// 	postTab = append(postTab, idata.username)
// 	postTab = append(postTab, idata.newPost)
// 	println(postTab[1])

// 	data := "test"

// 	err := tpl.Execute(w, data)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

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

type databaseInfo struct {
	idx     int
	pseudo  string
	message string
}
type comment struct {
	id      int
	image   string
	pseudo  string
	titre   string
	message string
}

type oui struct {
	Postee []comment
}

func database(username string, newPost string) {

	database, _ :=
		sql.Open("sqlite3", "test.db")
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
		// test = append(test, username)
		// test = append(test, newPost)
	}
	// fmt.Println(strconv.Itoa(id) + ": " + username + " " + newPost)
	// fmt.Println(test)
	// return test

}

func getInfo() []comment {

	database, _ :=
		sql.Open("sqlite3", "test.db")
	rows, _ :=
		database.Query("SELECT id, username, newPost FROM post")

	var _id int
	var test []comment
	var _pseudo string
	var _message string
	for rows.Next() {
		rows.Scan(&_id, &_pseudo, &_message)
		data := comment{
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
	newPost := r.FormValue("newPost")

	if userName != "" && newPost != "" {
		fmt.Println("tu rentre de dans wsh")
		database(userName, newPost)
		fmt.Println("tu sors de  wsh")
		// http.Redirect(w, r, "/main", http.StatusSeeOther)
	}

	// var infoPost []string
	// infoPost = append(infoPost, userName)
	// infoPost = append(infoPost, newPost)
	// println(&infoPost)

	tpl := template.Must(template.ParseFiles("assets/test.html"))

	data := oui{
		Postee: getInfo(),
	}

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

	http.HandleFunc("/post", RegisterHandle)
	http.ListenAndServe(":8082", nil)
}
