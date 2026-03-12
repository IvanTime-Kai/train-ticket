package grpc

import (
	"context"
	"fmt"

	proto "github.com/IvanTime-Kai/train-ticket-proto/gen/train"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TrainClient struct {
	client proto.TrainServiceClient
}

func NewTrainClient(host string, port int) (*TrainClient, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to train-server: %w", err)
	}

	return &TrainClient{
		client: proto.NewTrainServiceClient(conn),
	}, nil
}

func (c *TrainClient) ValidateSeats(ctx context.Context, tripID string, seatIDs []string) ([]*proto.SeatInfo, error) {
	resp, err := c.client.ValidateSeats(ctx, &proto.ValidateSeatsRequest{
		TripId:  tripID,
		SeatIds: seatIDs,
	})

	if err != nil {
		return nil, fmt.Errorf("grpc validate seats failed: %w", err)
	}

	if !resp.Valid {
		return nil, fmt.Errorf(resp.Message)
	}

	return resp.Seats, nil
}
