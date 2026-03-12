package grpc

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	proto "github.com/IvanTime-Kai/train-ticket-proto/gen/train"
	db "github.com/leminhthai/train-ticket/train-service/db/generated"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TrainGRPCServer struct {
	proto.UnimplementedTrainServiceServer
	q *db.Queries
}

func NewTrainGRPCServer(q *db.Queries) *TrainGRPCServer {
	return &TrainGRPCServer{
		q: q,
	}
}

func (s *TrainGRPCServer) ValidateSeats(ctx context.Context, req *proto.ValidateSeatsRequest) (*proto.ValidateSeatsResponse, error) {

	if req.TripId == "" {
		return nil, status.Error(codes.InvalidArgument, "trip_id is required")
	}

	if len(req.SeatIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "seat_ids is required")
	}

	trip, err := s.q.GetTripByID(ctx, req.TripId)
	if err != nil {
		if err == sql.ErrNoRows {
			return &proto.ValidateSeatsResponse{
				Valid:   false,
				Message: "trip not found",
			}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	seats := make([]*proto.SeatInfo, 0, len(req.SeatIds))
	for _, seatID := range req.SeatIds {
		seat, err := s.q.GetSeatByID(ctx, seatID)
		if err != nil {
			if err == sql.ErrNoRows {
				return &proto.ValidateSeatsResponse{
					Valid:   false,
					Message: fmt.Sprintf("seat %s not found", seatID),
				}, nil
			}
			return nil, status.Error(codes.Internal, err.Error())
		}

		// Verify seat thuộc train của trip
		if seat.TrainID != trip.TrainID {
			return &proto.ValidateSeatsResponse{
				Valid:   false,
				Message: fmt.Sprintf("seat %s does not belong to this trip", seatID),
			}, nil
		}

		price, _ := strconv.ParseFloat(seat.Price, 64)
		seats = append(seats, &proto.SeatInfo{
			SeatId:     seat.ID,
			SeatNumber: seat.SeatNumber,
			Class:      seat.Class,
			Price:      price,
		})
	}

	return &proto.ValidateSeatsResponse{
		Valid:   true,
		Message: "ok",
		Seats:   seats,
	}, nil
}
