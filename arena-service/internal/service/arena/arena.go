package arena

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/reymooy27/arena-backend/arena-service/db"
	pb "github.com/reymooy27/arena-backend/arena-service/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedArenaServiceServer
}

func (s *Server) GetArenas(ctx context.Context, in *emptypb.Empty) (*pb.ListArenaResponse, error) {
	query := `SELECT * FROM "arenas"`

	row, err := db.DB.QueryContext(ctx, query)

	if err != nil {
		slog.Error("Query arena", "err", err)
		return nil, fmt.Errorf("Cannot get arena data")
	}

	defer row.Close()

	var arenas []*pb.ArenaData
	var createdAt time.Time

	for row.Next() {
		var arena pb.ArenaData
		if err := row.Scan(&arena.ArenaId, &createdAt, &arena.Name, &arena.Description); err != nil {
			slog.Error("Scan arena", "err", err)
			return nil, fmt.Errorf("Cannot get arena data")
		}

		arena.CreatedAt = timestamppb.New(createdAt)
		arenas = append(arenas, &arena)
	}

	return &pb.ListArenaResponse{Arenas: arenas}, nil
}

func (s *Server) GetArenaById(ctx context.Context, in *pb.GetArenaRequest) (*pb.ArenaData, error) {
	id := in.ArenaId

	if id == 0 {
		return nil, fmt.Errorf("Invalid arena id")
	}

	query := `SELECT * FROM "arenas" WHERE "id" = $1`

	row := db.DB.QueryRowContext(ctx, query, id)

	var arena pb.ArenaData
	var createdAt time.Time

	if err := row.Scan(&arena.ArenaId, &createdAt, &arena.Name, &arena.Description); err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Scan arena", "err", err)
			return nil, fmt.Errorf("Cannot get arena data")
		}

		return nil, fmt.Errorf("Cannot get arena data")
	}
	arena.CreatedAt = timestamppb.New(createdAt)

	return &arena, nil
}

func (s *Server) CreateArena(ctx context.Context, in *pb.ArenaRequest) (*pb.ArenaResponse, error) {

	query := `INSERT INTO "arenas" ("name","description") VALUES ($1, $2)`
	_, err := db.DB.ExecContext(ctx, query, in.Name, in.Description)

	if err != nil {
		slog.Error("Query insert arena", "err", err)
		return &pb.ArenaResponse{Message: "Failed create arena", Success: false}, fmt.Errorf("Cannot create arena")
	}

	return &pb.ArenaResponse{Message: "Successfully create arena", Success: true}, nil
}

func (s *Server) DeleteArena(ctx context.Context, in *pb.GetArenaRequest) (*pb.ArenaResponse, error) {
	id := in.ArenaId

	if id == 0 {
		return nil, fmt.Errorf("Invalid arena id")
	}

	query := `DELETE FROM "arenas" WHERE "id" = $1`
	row, err := db.DB.ExecContext(ctx, query, id)

	rowsAffected, err := row.RowsAffected()

	if err != nil {
		slog.Error("Cannot delete arena", "err", err)
		return &pb.ArenaResponse{Message: "Failed delete arena", Success: false}, fmt.Errorf("Cannot delete arena")
	}

	if rowsAffected == 0 {
		return &pb.ArenaResponse{Message: "No data with the id", Success: false}, fmt.Errorf("Cannot delete arena")
	}

	return &pb.ArenaResponse{Message: "Successfully delete arena", Success: true}, nil
}

func (s *Server) UpdateArena(ctx context.Context, in *pb.UpdateArenaRequest) (*pb.ArenaResponse, error) {
	id := in.ArenaId

	if id == 0 {
		return nil, fmt.Errorf("Invalid arena id")
	}

	query := `UPDATE "arenas" SET name = $1, description = $2 WHERE "id" = $3`
	result, err := db.DB.ExecContext(ctx, query, in.Name, in.Description, id)

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		slog.Error("Cannot update arena", "err", err)
		return &pb.ArenaResponse{Message: "Failed delete arena", Success: false}, fmt.Errorf("Cannot update arena")
	}

	if rowsAffected == 0 {
		return &pb.ArenaResponse{Message: "No data with the id", Success: false}, fmt.Errorf("No data with the id")
	}

	return &pb.ArenaResponse{Message: "Successfully update data", Success: true}, nil
}
