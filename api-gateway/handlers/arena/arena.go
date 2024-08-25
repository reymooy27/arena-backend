package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/reymooy27/arena-backend/api-gateway/handlers/gateway"
	pb "github.com/reymooy27/arena-backend/api-gateway/proto/arena"
	"github.com/reymooy27/arena-backend/api-gateway/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Request struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type ArenaResponse struct {
	ArenaId     int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type ArenaHandler struct {
	arenaClient pb.ArenaServiceClient
}

func NewArenaHandler(apiGateway gateway.APIGateway) *ArenaHandler {
	return &ArenaHandler{arenaClient: apiGateway.ArenaClient}
}

func (s *ArenaHandler) CreateArena(w http.ResponseWriter, r *http.Request) {
	var data Request

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		slog.Error("Cannot decode json", "err", err)
		utils.JSONResponse(w, 400, "Invalid request body")
		return
	}

	req := &pb.ArenaRequest{
		Name:        data.Name,
		Description: data.Description,
	}

	res, err := s.arenaClient.CreateArena(context.Background(), req)
	if err != nil {
		slog.Error("Cannot create arena", "err", err)
		utils.JSONResponse(w, 500, "Cannot create arena")
		return
	}

	utils.JSONResponse(w, 200, res)
}

func (s *ArenaHandler) GetArenaById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("Invalid arena id", "err", err)
		utils.JSONResponse(w, 500, "Invalid arena id")
		return
	}

	req := &pb.GetArenaRequest{
		ArenaId: parsedId,
	}

	res, err := s.arenaClient.GetArenaById(context.Background(), req)
	if err != nil {
		slog.Error("Cannot get arena data", "err", err)
		utils.JSONResponse(w, 500, "Cannot get arena data")
		return
	}
	response := ArenaResponse{
		ArenaId:     res.ArenaId,
		Name:        res.Name,
		Description: res.Description,
		CreatedAt:   res.CreatedAt.AsTime(),
	}

	utils.JSONResponse(w, 200, response)
}

func (s *ArenaHandler) GetArenas(w http.ResponseWriter, r *http.Request) {

	res, err := s.arenaClient.GetArenas(context.Background(), &emptypb.Empty{})
	if err != nil {
		slog.Error("Cannot get arena data", "err", err)
		utils.JSONResponse(w, 500, "Cannot get arena data")
		return
	}

	response := make([]ArenaResponse, 0, len(res.Arenas))

	for _, v := range res.Arenas {
		response = append(response,
			ArenaResponse{
				ArenaId:     v.ArenaId,
				CreatedAt:   v.CreatedAt.AsTime(),
				Name:        v.Name,
				Description: v.Description,
			},
		)

	}

	utils.JSONResponse(w, 200, response)
}

func (s *ArenaHandler) UpdateArena(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("Invalid arena id", "err", err)
		utils.JSONResponse(w, 500, "Invalid arena id")
		return
	}

	var data Request

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		slog.Error("Cannot decode json", "err", err)
		utils.JSONResponse(w, 400, "Invalid request body")
		return
	}

	req := &pb.UpdateArenaRequest{
		ArenaId:     parsedId,
		Name:        data.Name,
		Description: data.Description,
	}

	res, err := s.arenaClient.UpdateArena(context.Background(), req)
	if err != nil {
		slog.Error("Cannot update arena data", "err", err)
		utils.JSONResponse(w, 500, "Cannot update arena data")
		return
	}

	utils.JSONResponse(w, 200, res)
}

func (s *ArenaHandler) DeleteArena(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("Invalid arena id", "err", err)
		utils.JSONResponse(w, 500, "Invalid arena id")
		return
	}

	req := &pb.GetArenaRequest{
		ArenaId: parsedId,
	}

	res, err := s.arenaClient.DeleteArena(context.Background(), req)
	if err != nil {
		slog.Error("Cannot delete arena data", "err", err)
		utils.JSONResponse(w, 500, "Cannot delete arena data")
		return
	}

	utils.JSONResponse(w, 200, res)
}
