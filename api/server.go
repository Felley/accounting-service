package api

import (
	"context"

	"github.com/Felley/accounting-service/protos/accounting"
)

// RateRequest ...
type RateRequest struct {
	accounting.UnimplementedCurrencyServer
}

// GetRate ...
func (api *RateRequest) GetRate(ctx context.Context, rr *accounting.RateRequest) (*accounting.RateResponce, error) {
	return &accounting.RateResponce{Rate: 1.0}, nil
}

func (api *RateRequest) mustEmbedUnimplementedCurrencyServer() {

}
