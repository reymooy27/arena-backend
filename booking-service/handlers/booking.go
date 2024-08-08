package handlers

import (
	// "database/sql"
	"encoding/json"
	"fmt"
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
	Username     string    `json:"username"`
	ArenaName    string    `json:"arena_name"`
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

	query := `SELECT * FROM bookings WHERE user_id = $1 ORDER BY created_at DESC LIMIT 5 `
	rows, err := db.DB.Query(query, user.Id)
	if err != nil {
		slog.Error("Query error", "message", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Cannot get data")
		return
	}

	defer rows.Close()

	var bookings []Booking
	userDataChan := make(chan *UserDataResponse)
	userErrChan := make(chan error)

	go func() {
		userData, err := GetUserData(user.Id)
		if err != nil {
			userErrChan <- err
		}
		userDataChan <- userData
	}()

	var userData *UserDataResponse
	select {
	case userData = <-userDataChan:
	case userErr := <-userErrChan:
		slog.Error("Error get user data", "message", userErr)
		utils.JSONResponse(w, http.StatusInternalServerError, "Cannot get data")
		return
	}
	if userData == nil {
		slog.Error("userData is nil", "userId", user.Id)
		utils.JSONResponse(w, http.StatusInternalServerError, "Cannot get data")
		return
	}

	for rows.Next() {
		var booking Booking
		if err := rows.Scan(&booking.Id, &booking.ArenaId, &booking.UserId, &booking.CreatedAt, &booking.BookingSlots); err != nil {
			slog.Error("Scan error", "message", err)
			utils.JSONResponse(w, http.StatusInternalServerError, "Cannot get data")
			return
		}

		arenaDataChan := make(chan *ArenaDataResponse)
		arenaErrChan := make(chan error)

		go func(arenaId int) {
			arenaData, err := GetArenaData(arenaId)
			if err != nil {
				arenaErrChan <- err
			}
			arenaDataChan <- arenaData
		}(booking.ArenaId)

		var arenaData *ArenaDataResponse
		select {
		case arenaData = <-arenaDataChan:
		case arenaErr := <-arenaErrChan:
			slog.Error("Error get arena data", "message", arenaErr)
			continue
		}
		if arenaData == nil {
			slog.Error("arenaData is nil", "arenaId", booking.ArenaId)
			utils.JSONResponse(w, http.StatusInternalServerError, "Cannot get data")
			return
		}

		booking.Username = userData.Username
		booking.ArenaName = arenaData.Nama
		bookings = append(bookings, booking)
	}

	utils.JSONResponse(w, http.StatusOK, bookings)
}

type UserDataResponse struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
}

func GetUserData(userId int) (*UserDataResponse, error) {
	res, err := http.Get(fmt.Sprintf("http://localhost:8001/user/%d", userId))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, err
	}

	var response UserDataResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

type ArenaDataResponse struct {
	Nama string `json:"nama"`
	Id   int    `json:"id"`
}

func GetArenaData(arenaId int) (*ArenaDataResponse, error) {
	res, err := http.Get(fmt.Sprintf("http://localhost:8000/arena/%d", arenaId))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, err
	}
	var response ArenaDataResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
