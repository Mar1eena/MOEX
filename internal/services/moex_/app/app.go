package app

import (
	"context"
	"net"
	"os"

	"github.com/Mar1eena/Test_gRPC/internal/pkg/zlog"
	"github.com/Mar1eena/Test_gRPC/internal/services/moex_/server"
	"google.golang.org/grpc"
)

func App() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zl := zlog.New()

	// grpc server
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}
	gs := grpc.NewServer()

	// listener
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		zl.Fatal("failed to listen", err)
	}
	zl.Info("server listening at %v", lis.Addr().String())

	// service
	service := server.Service(zl)

	// register service
	server.RegisterMOEXServer(gs, service)

	go func() {
		<-ctx.Done()
		gs.GracefulStop()
	}()

	go func() {
		if err := gs.Serve(lis); err != nil {
			zl.Error("failed to serve", err)
		}
	}()

	// Wait for signal
	<-ctx.Done()
}
