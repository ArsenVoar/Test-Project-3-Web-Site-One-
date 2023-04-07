package main

import (
	"fmt"
	"net/http"
	"text/template"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Comments struct {
	Id                      uint16
	Title, Anons, Full_Text string
}

var posts = []Comments{}

func mainPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/mainPage.html", "templates/header.html", "templates/footer.html")

	t.ExecuteTemplate(w, "mainPage", nil)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")

	t.ExecuteTemplate(w, "create", nil)
}

func examples(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/examples.html", "templates/header.html", "templates/footer.html")

	t.ExecuteTemplate(w, "examples", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "No")
	} else {
		db, err := sql.Open("mysql", "root@tcp(localhost:3306)/test-project")

		if err != nil {
			panic(err)
		}

		defer db.Close()

		inst, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES('%s', '%s', '%s')", title, anons, full_text))
		if err != nil {
			panic(err)
		}
		defer inst.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func comments(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/comments.html", "templates/header.html", "templates/footer.html")

	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/test-project")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	res, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		panic(err)
	}

	posts = []Comments{}
	for res.Next() {
		var post Comments
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_Text)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}
	t.ExecuteTemplate(w, "comments", posts)
}

func HandleFunc() {
	router := mux.NewRouter()
	router.HandleFunc("/", mainPage)
	router.HandleFunc("/create", create)
	router.HandleFunc("/examples", examples)
	router.HandleFunc("/comments", comments)
	router.HandleFunc("/save_article", save_article)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8888", nil)
}

func main() {
	HandleFunc()
}
