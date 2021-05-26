package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func login(logUser string, logPassword string) {
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
			break
		}
	}
	database.Close()
}

func main() {
	login("PetitCul_Pâle", "123456789")
}
