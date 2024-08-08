package handlers

import (
	// "database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	// "os"
	// "strings"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/reymooy27/arena-backend/booking-service/db"
	"github.com/reymooy27/arena-backend/booking-service/utils"
)

type Body struct {
	ArenaId int       `json:"arena_id"`
	Time    time.Time `json:"time"`
}

type Booking struct {
	Id           int       `json:"id"`
	ArenaId      int       `json:"arena_id"`
	UserId       int       `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	BookingSlots string    `json:"booking_slots"`
}

type Claim struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// INFO: not finished,
// INFO: still testing
func CreateBooking(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(*Claim)

	var body Body

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	query := `INSERT INTO bookings (arena_id, user_id, booking_slots) VALUES ($1, $2, $3)`
	result, err := db.DB.Exec(query, body.ArenaId, user.Id, "slot 1")
	log.Println(result)
	if err != nil {
		slog.Error("Query error", "message", err)
		utils.JSONResponse(w, http.StatusBadRequest, "Cannot create booking")
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Successfully create booking")
}

// TODO: join table/data from different database (microrservices)
func GetUserBookings(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(*Claim)

	query := `SELECT * FROM bookings WHERE user_id = $1`
	rows, err := db.DB.Query(query, user.Id)
	if err != nil {
		slog.Error("Query error", "message", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Cannot get data")
		return
	}

	defer rows.Close()

	var bookings []Booking

	for rows.Next() {
		var booking Booking
		if err := rows.Scan(&booking.Id, &booking.ArenaId, &booking.UserId, &booking.CreatedAt, &booking.BookingSlots); err != nil {
			slog.Error("Scan error", "message", err)
			utils.JSONResponse(w, http.StatusInternalServerError, "Cannot get data")
			return
		}
		bookings = append(bookings, booking)
	}

	utils.JSONResponse(w, http.StatusOK, bookings)
}
