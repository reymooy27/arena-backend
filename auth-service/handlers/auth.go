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

	query := `SELECT * FROM "user" WHERE username = $1`
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
		slog.Info("Invalid request body", err)
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

	var existedUsername string
	existedUsernameQuery := `
  SELECT username 
  FROM "user" 
  WHERE username = $1 
  LIMIT 1`
	err = db.DB.QueryRow(existedUsernameQuery, strings.Trim(data.Username, " ")).Scan(&existedUsername)
	if err != nil && err != sql.ErrNoRows {
		slog.Error("Error", err)
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

	query := `INSERT INTO "user" ("username","password") VALUES ($1, $2)`
	result, err := db.DB.Exec(query, data.Username, string(hashedPassword))
	log.Println(result)

	if err != nil {
		slog.Error("Error", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Cannot create user")
		return
	}

	utils.JSONResponse(w, http.StatusCreated, "Successfully create account")
}

func Verify(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	utils.JSONResponse(w, http.StatusOK, user)
}
