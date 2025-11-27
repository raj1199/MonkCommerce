package repo

import (
	"context"

	"github.com/you/monk-coupons/pkg/model"
)

type Repository interface {
	Create(ctx context.Context, c *model.Coupon) (string, error)
	GetAll(ctx context.Context) ([]*model.Coupon, error)
	GetByID(ctx context.Context, id string) (*model.Coupon, error)
	Update(ctx context.Context, id string, c *model.Coupon) error
	Delete(ctx context.Context, id string) error
}
