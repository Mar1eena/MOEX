package server

import (
	"context"

	moex_contract_v1 "github.com/Mar1eena/trb_proto/gen/go/moex"
)

func (s *MOEXServiceServer) Logon(ctx context.Context, req *moex_contract_v1.Logonrequest) (*moex_contract_v1.Logonresponse, error) {
	resp, err := Logon(req)
	if err != nil {
		return nil, err
	}
	return &moex_contract_v1.Logonresponse{Reply: resp}, nil
}
