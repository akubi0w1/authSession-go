package handler

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"auth-session/server"
	"github.com/google/uuid"
)

// Signup signup
func Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("signup")
		// read body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			server.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var req signupRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			server.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// TODO: データのバリデーション

		// password encrypt
		hash, err := server.PasswordHash(req.Password)
		if err != nil {
			server.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// insert data
		_, err = db.Exec(
			"INSERT INTO users(id, name, password) VALUES (?,?,?)",
			req.ID,
			req.Name,
			hash,
		)

		// response
		server.Success(w, &signupResponse{
			ID:       req.ID,
			Name:     req.Name,
			Password: hash,
		})
	}
}

type signupRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type signupResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Login login
func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("login")
		// read body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			server.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var req loginRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			server.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// auth
		var hash string
		row := db.QueryRow("SELECT password FROM users WHERE id=?", req.ID)
		if err = row.Scan(&hash); err != nil {
			server.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		err = server.PasswordVerify(hash, req.Password)
		if err != nil {
			server.ErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// create session id
		sessionID, err := uuid.NewRandom()
		if err != nil {
			server.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// update session
		_, err = db.Exec(
			"UPDATE users SET session_id = ? WHERE id=?",
			sessionID.String(),
			req.ID,
		)
		if err != nil {
			server.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// create cookie
		cookie := http.Cookie{
			Name:  "_cookie",
			Value: sessionID.String(),
		}
		http.SetCookie(w, &cookie)

		// response
		server.Success(w, &messageResponse{Message: "success"})
	}
}

type loginRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

// Logout logout
func Logout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("logout")
		// check cookie
		cookie, err := r.Cookie("_cookie")
		if err != nil {
			server.ErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// cookieの無効化
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)

		server.Success(w, &messageResponse{Message: "success"})
	}
}

type messageResponse struct {
	Message string `json:"message"`
}
