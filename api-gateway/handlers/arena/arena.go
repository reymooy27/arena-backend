package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/reymooy27/arena-backend/api-gateway/handlers/gateway"
	pb "github.com/reymooy27/arena-backend/api-gateway/proto/arena"
	"github.com/reymooy27/arena-backend/api-gateway/utils"
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
	ArenaId     int64
	Name        string
	Description string
	CreatedAt   string
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
		utils.JSONResponse(w, 400, fmt.Errorf("Invalid request body"))
		return
	}

	req := &pb.ArenaRequest{
		Name:        data.Name,
		Description: data.Description,
	}

	res, err := s.arenaClient.CreateArena(context.Background(), req)
	if err != nil {
		slog.Error("Cannot create arena", "err", err)
		utils.JSONResponse(w, 500, fmt.Errorf("Cannot create arena"))
		return
	}

	var response = Response{
		Message: res.Message,
		Success: res.Success,
	}

	utils.JSONResponse(w, 200, response)
}

func (s *ArenaHandler) GetArenaById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("Invalid arena id", "err", err)
		utils.JSONResponse(w, 500, fmt.Errorf("Invalid arena id"))
		return
	}

	req := &pb.GetArenaRequest{
		ArenaId: parsedId,
	}

	res, err := s.arenaClient.GetArenaById(context.Background(), req)
	if err != nil {
		slog.Error("Cannot get arena data", "err", err)
		utils.JSONResponse(w, 500, fmt.Errorf("Cannot get arena data"))
		return
	}

	var response = ArenaResponse{
		ArenaId:     res.ArenaId,
		Name:        res.Name,
		Description: res.Description,
	}

	utils.JSONResponse(w, 200, response)
}
