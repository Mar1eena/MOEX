package server

import (
	"context"

	moex_contract_v1 "github.com/Mar1eena/trb_proto/gen/go/moex"
)

func (s *MOEXServiceServer) Dealing(ctx context.Context, req *moex_contract_v1.DealingRequest) (*moex_contract_v1.DealingResponse, error) {
	return Dealing(req)
}
