package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/reymooy27/arena-backend/auth-service/db"
	"github.com/reymooy27/arena-backend/auth-service/utils"
)

type Data struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
}

type Claim struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {

	var data Data

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if data.Username == "" || data.Password == "" {
		utils.JSONResponse(w, http.StatusBadRequest, "Username and password cannot be empty")
		return
	}

	var user User

	query := `SELECT * FROM "users" WHERE username = $1`
	row := db.DB.QueryRow(query, data.Username)
	err = row.Scan(&user.Id, &user.CreatedAt, &user.Username, &user.Password)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Wrong username")
		slog.Error("Error", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))

	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Wrong password")
		slog.Error("Cannot compare hash", err)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := Claim{
		Id:       user.Id,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Println("SECRET not set")
	}

	jwt, err := token.SignedString([]byte(secret))
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Cannot create token")
		slog.Error("Cannot sign jwt", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    jwt,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	})

	utils.JSONResponse(w, http.StatusOK, "Successfully login")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	})

	utils.JSONResponse(w, http.StatusOK, "Logout successful")
}

func Signup(w http.ResponseWriter, r *http.Request) {

	var data Data

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if data.Username == "" || data.Password == "" {
		utils.JSONResponse(w, http.StatusBadRequest, "Username and password cannot be empty")
		return
	}

	if len(data.Password) < 8 {
		utils.JSONResponse(w, http.StatusBadRequest, "Password should be at least 8 characters")
		return
	}

	if len(data.Password) > 64 {
		utils.JSONResponse(w, http.StatusBadRequest, "Password should be at most 64 characters")
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		slog.Error("Cannot begin transaction", "message", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				log.Fatalf("Failed to commit transaction: %v", err)
			}
		}
	}()

	var existedUsername string
	existedUsernameQuery := ` SELECT username FROM "users" WHERE username = $1 LIMIT 1`
	err = tx.QueryRow(existedUsernameQuery, strings.Trim(data.Username, " ")).Scan(&existedUsername)
	if err != nil && err != sql.ErrNoRows {
		slog.Error("Query check username exist", "message", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	if existedUsername != "" {
		utils.JSONResponse(w, http.StatusInternalServerError, "Username already exist")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Cannot hashed Password")
		return
	}

	var userId int
	query := `INSERT INTO "users" ("username","password") VALUES ($1, $2) RETURNING "id"`
	if err = tx.QueryRow(query, data.Username, string(hashedPassword)).Scan(&userId); err != nil {
		slog.Error("Query create user", "message", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Cannot create user")
		return
	}

	slog.Info("id", userId)

	query = `INSERT INTO "profiles" ("user_id") VALUES ($1)`
	_, err = tx.Exec(query, userId)
	if err != nil {
		slog.Error("Query create profile", "message", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Cannot create user")
		return
	}

	utils.JSONResponse(w, http.StatusCreated, "Successfully create account")
}

func Verify(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	utils.JSONResponse(w, http.StatusOK, user)
}

// TODO : Need to check the previous password ???
func ChangePassword(w http.ResponseWriter, r *http.Request) {

	context := r.Context().Value("user").(*Claim)

	type ChangePasswordBody struct {
		NewPassword        string `json:"new_password"`
		ConfirmNewPassword string `json:"confirm_password"`
	}

	var body ChangePasswordBody
	var user User

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid body request")
		return
	}

	query := `SELECT id, username FROM "users" WHERE id = $1 LIMIT 1`

	if err := db.DB.QueryRow(query, context.Id).Scan(&user.Id, &user.Username); err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Query error", err)
			utils.JSONResponse(w, http.StatusNotFound, "Cannot find user")
			return
		}

		slog.Error("Query error", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "There is something wrong")
		return
	}

	if user.Id != context.Id {
		slog.Error("Forbidden change password")
		utils.JSONResponse(w, http.StatusForbidden, "Cannot change password")
		return
	}

	if body.NewPassword != body.ConfirmNewPassword {
		utils.JSONResponse(w, http.StatusBadRequest, "New password must match confirm password")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Cannot hashed Password")
		return
	}

	query = `UPDATE "users" SET password = $1 WHERE id = $2`
	_, err = db.DB.Exec(query, string(hashedPassword), user.Id)

	if err != nil {
		slog.Error("Error", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Cannot change password")
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Successfuly change password")
}
