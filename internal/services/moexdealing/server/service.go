package server

import (
	"context"

	moex_contract_v1 "github.com/Mar1eena/trb_proto/gen/go/moex"
)

func (s *MOEXServiceServer) Request(ctx context.Context, req *moex_contract_v1.Req) (*moex_contract_v1.Resp, error) {
	return &moex_contract_v1.Resp{Reply: ""}, nil
}

func (s *MOEXServiceServer) Dealing(ctx context.Context, req *moex_contract_v1.DealingRequest) (*moex_contract_v1.DealingResponse, error) {
	resp, err := Dealing(req)
	return &moex_contract_v1.DealingResponse{Reply: resp}, err
}
