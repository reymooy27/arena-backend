package handlers

import (
	// "database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

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

type ArenaDataResponse struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type UserDataResponse struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
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

// TODO: still looking for better aproach to join data between services
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
		booking.ArenaName = arenaData.Name
		bookings = append(bookings, booking)
	}

	utils.JSONResponse(w, http.StatusOK, bookings)
}

func GetUserData(userId int) (*UserDataResponse, error) {
	userServiceURL := os.Getenv("USER_SERVICE_URL")

	url := fmt.Sprintf("%s/user/%d", userServiceURL, userId)
	res, err := http.Get(url)
	if err != nil {
		slog.Error("Cannot fetch user data", "err", err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.Error("Cannot fetch user data", "err", err)
		return nil, err
	}

	var response UserDataResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func GetArenaData(arenaId int) (*ArenaDataResponse, error) {

	arenaServiceURL := os.Getenv("ARENA_SERVICE_URL")

	res, err := http.Get(fmt.Sprintf("%s/arena/%d", arenaServiceURL, arenaId))
	if err != nil {
		slog.Error("Cannot fetch arena data", "err", err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.Error("Cannot fetch arena data", "err", err)
		return nil, err
	}
	var response ArenaDataResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
