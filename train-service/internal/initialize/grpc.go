package initialize

import (
	"fmt"
	"net"

	proto "github.com/IvanTime-Kai/train-ticket-proto/gen/train"
	db "github.com/leminhthai/train-ticket/train-service/db/generated"
	"github.com/leminhthai/train-ticket/train-service/global"
	grpcServer "github.com/leminhthai/train-ticket/train-service/internal/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitGRPC() {
	queries := db.New(global.Mdb)

	srv := grpcServer.NewTrainGRPCServer(queries)
	grpcPort := global.Config.GRPC.Port

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))

	if err != nil {
		global.Logger.Fatal("failed to listen gRPC", zap.Error(err))
	}

	s := grpc.NewServer()
	proto.RegisterTrainServiceServer(s, srv)

	global.Logger.Info("gRPC server started", zap.Int("port", grpcPort))

	// Run gRPC in goroutine - not block HTTP server
	go func() {
		if err := s.Serve(lis); err != nil {
			global.Logger.Fatal("gRPC server failed", zap.Error(err))
		}
	}()
}
