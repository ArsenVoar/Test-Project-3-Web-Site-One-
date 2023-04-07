package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/test-project")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	inst, err := db.Query("INSERT INTO `comments`(`Name`, `NickName`) VALUES('Arsen', 'VoAr')")
	if err != nil {
		panic(err)
	}
	defer inst.Close()

	fmt.Println("Connected")
}
