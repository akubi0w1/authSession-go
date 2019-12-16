package main

import (
	"auth-session/handler"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// connect db
	DB, err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	// routing
	http.HandleFunc("/login", handler.Login(DB))
	http.HandleFunc("/logout", handler.Logout(DB))
	http.HandleFunc("/signup", handler.Signup(DB))

	http.HandleFunc("/whoami", handler.WhoAmI(DB))

	// start server
	log.Println("srever running...")
	http.ListenAndServe(":8080", nil)
}

func initDB() (*sql.DB, error) {
	return sql.Open("mysql", "root:password@tcp(localhost:3307)/auth_sample")
}
