package handler

import (
	"auth-session/server"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

// WhoAmI get user name by cookie
func WhoAmI(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("who am i")
		// TODO: セッションチェック！
		cookie, err := r.Cookie("_cookie")
		if err != nil {
			server.ErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// TODO: userIDの取り出し
		var name string
		row := db.QueryRow("SELECT name FROM users WHERE session_id=?", cookie.Value)
		if err = row.Scan(&name); err != nil {
			server.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// TODO: who am i!
		server.Success(w, &whoAmIResponse{
			Message: fmt.Sprintf("I am %s", name),
		})
	}
}

type whoAmIResponse struct {
	Message string `json:"message"`
}
