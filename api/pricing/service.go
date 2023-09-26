package pricing

import (
	"context"
	"time"
)

type Service struct{}

func (s *Service) CalculatePricing(ctx context.Context, dateParam string) (
	CalculatePricingResp, error) {
	var (
		resp = CalculatePricingResp{}
	)
	date, err := time.Parse("02012006", dateParam)
	if err != nil {
		resp.ErrCode = 40
		resp.ErrMessage = "date parameter invalid"
		return resp, err
	}
	resp.Price = date.Day() * 100000
	return resp, nil
}

func NewService() *Service {
	return &Service{}
}
