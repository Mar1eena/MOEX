package server

import (
	"github.com/Mar1eena/Test_gRPC/internal/pkg/zlog"
	moex_contract_v1 "github.com/Mar1eena/trb_proto/gen/go/moex"
	"google.golang.org/grpc"
)

// Init service
type MOEXServiceServer struct {
	moex_contract_v1.UnimplementedMoexServer
	zl zlog.ZLogger
}

func RegisterMOEXServer(srv *grpc.Server, service *MOEXServiceServer) {
	moex_contract_v1.RegisterMoexServer(srv, service)
}

// Init nats
func Service(zl zlog.ZLogger) *MOEXServiceServer {
	return &MOEXServiceServer{
		zl: zl,
	}
}
