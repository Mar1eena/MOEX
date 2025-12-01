package app

import (
	"context"
	"net"
	"os"

	"github.com/Mar1eena/Test_gRPC/internal/pkg/zlog"
	"github.com/Mar1eena/Test_gRPC/internal/services/moexdealing/server"
	"google.golang.org/grpc"
)

func App() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zl := zlog.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}
	gs := grpc.NewServer()

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		zl.Fatal("failed to listen", err)
	}
	zl.Info("server listening at %v", lis.Addr().String())

	service := server.Service(zl)

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

	<-ctx.Done()
}
