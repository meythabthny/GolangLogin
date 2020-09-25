package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

import "golang.org/x/crypto/bcrypt"

import "net/http"

import "fmt"

var db *sql.DB
var err error

func daftarPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "daftar.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	namalengkap := req.FormValue("namalengkap")
	var user string

	err := db.QueryRow("SELECT username FROM user WHERE username=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Pendaftaran Gagal.", 500)
			return
		}

		_, err = db.Exec("INSERT INTO user(username, password,namalengkap) VALUES(?, ?, ?)", username, hashedPassword,namalengkap)
		if err != nil {
			http.Error(res, "Pendaftaran Gagal.", 500)
			return
		}

		res.Write([]byte("Pendaftaran Berhasil"))
		return
	case err != nil:
		http.Error(res, "Pendaftaran Gagal", 500)
		return
	default:
		http.Redirect(res, req, "/", 301)
	}
}

func loginPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "login.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var databaseUsername string
	var databasePassword string

	err := db.QueryRow("SELECT username, password FROM user WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	res.Write([]byte("Hello" + " " + databaseUsername))

}

func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "home.html")
}

func main() {
	db, err = sql.Open("mysql", "root:@/cobagolang")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}



	http.HandleFunc("/daftar", daftarPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/", homePage)
	var address = "localhost:9000"
    	fmt.Printf("server started at %s\n", address)
	http.ListenAndServe(":9000", nil)
}
