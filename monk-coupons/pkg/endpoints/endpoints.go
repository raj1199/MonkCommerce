package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/you/monk-coupons/pkg/model"
	"github.com/you/monk-coupons/pkg/service"
)

type Endpoints struct {
	Create      endpoint.Endpoint
	GetAll      endpoint.Endpoint
	GetByID     endpoint.Endpoint
	Update      endpoint.Endpoint
	Delete      endpoint.Endpoint
	Applicable  endpoint.Endpoint
	Apply       endpoint.Endpoint
}

func Make(s service.Service) Endpoints {
	return Endpoints{
		Create: func(ctx context.Context, req interface{}) (interface{}, error) {
			return s.CreateCoupon(ctx, req.(*model.Coupon))
		},
		GetAll: func(ctx context.Context, req interface{}) (interface{}, error) {
			return s.GetAll(ctx)
		},
		GetByID: func(ctx context.Context, req interface{}) (interface{}, error) {
			return s.GetByID(ctx, req.(string))
		},
		Update: func(ctx context.Context, req interface{}) (interface{}, error) {
			r := req.(map[string]interface{})
			return nil, s.Update(ctx, r["id"].(string), r["payload"].(*model.Coupon))
		},
		Delete: func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, s.Delete(ctx, req.(string))
		},
		Applicable: func(ctx context.Context, req interface{}) (interface{}, error) {
			return s.ApplicableCoupons(ctx, req.(*model.Cart))
		},
		Apply: func(ctx context.Context, req interface{}) (interface{}, error) {
			r := req.(map[string]interface{})
			return s.ApplyCoupon(ctx, r["id"].(string), r["cart"].(*model.Cart))
		},
	}
}
