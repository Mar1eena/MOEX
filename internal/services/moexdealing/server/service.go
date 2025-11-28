package server

import (
	"context"

	moex_contract_v1 "github.com/Mar1eena/trb_proto/gen/go/moex"
)

func (s *MOEXServiceServer) Request(ctx context.Context, req *moex_contract_v1.Req) (*moex_contract_v1.Resp, error) {
	resp, err := Request(req)
	if err != nil {
		return nil, err
	}
	return &moex_contract_v1.Resp{Reply: resp}, nil
}
