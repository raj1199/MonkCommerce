package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/you/monk-coupons/pkg/model"
	"github.com/you/monk-coupons/pkg/repo"
)

type Service interface {
	CreateCoupon(ctx context.Context, c *model.Coupon) (string, error)
	GetAll(ctx context.Context) ([]*model.Coupon, error)
	GetByID(ctx context.Context, id string) (*model.Coupon, error)
	Update(ctx context.Context, id string, c *model.Coupon) error
	Delete(ctx context.Context, id string) error
	ApplicableCoupons(ctx context.Context, cart *model.Cart) ([]model.ApplicableResponse, error)
	ApplyCoupon(ctx context.Context, id string, cart *model.Cart) (*model.Cart, error)
}

type service struct {
	r repo.Repository
}

func NewService(r repo.Repository) Service {
	return &service{r}
}

func (s *service) CreateCoupon(ctx context.Context, c *model.Coupon) (string, error) {
	if err := model.ValidateCoupon(c); err != nil {
		return "", err
	}
	return s.r.Create(ctx, c)
}

func (s *service) GetAll(ctx context.Context) ([]*model.Coupon, error) {
	return s.r.GetAll(ctx)
}

func (s *service) GetByID(ctx context.Context, id string) (*model.Coupon, error) {
	return s.r.GetByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id string, c *model.Coupon) error {
	if err := model.ValidateCoupon(c); err != nil {
		return err
	}
	return s.r.Update(ctx, id, c)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.r.Delete(ctx, id)
}

func (s *service) ApplicableCoupons(ctx context.Context, cart *model.Cart) ([]model.ApplicableResponse, error) {
	out := []model.ApplicableResponse{}

	cart.CalcTotal()
	all, err := s.r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, c := range all {
		if c.ExpiresAt != nil && time.Now().After(*c.ExpiresAt) {
			continue
		}

		discount := s.estimate(c, cart)
		if discount > 0 {
			out = append(out, model.ApplicableResponse{
				CouponID: c.ID,
				Type:     string(c.Type),
				Discount: discount,
			})
		}
	}
	return out, nil
}

func (s *service) ApplyCoupon(ctx context.Context, id string, cart *model.Cart) (*model.Cart, error) {
	c, err := s.r.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("coupon not found")
	}

	if c.ExpiresAt != nil && time.Now().After(*c.ExpiresAt) {
		return nil, fmt.Errorf("coupon expired")
	}

	cart.CalcTotal()

	switch c.Type {
	case model.CartWise:
		d := c.Details.(*model.CartWiseDetails)
		if cart.TotalPrice >= d.Threshold {
			discount := model.Round(cart.TotalPrice * d.Discount / 100)
			cart.TotalDiscount = discount
			cart.FinalPrice = model.Round(cart.TotalPrice - discount)
		}

	case model.ProductWise:
		d := c.Details.(*model.ProductWiseDetails)
		for i := range cart.Items {
			if cart.Items[i].ProductID == d.ProductID {
				itemDiscount := float64(cart.Items[i].Quantity) * cart.Items[i].Price * d.Discount / 100
				cart.Items[i].TotalDiscount = model.Round(itemDiscount)
				cart.TotalDiscount += cart.Items[i].TotalDiscount
			}
		}
		cart.FinalPrice = model.Round(cart.TotalPrice - cart.TotalDiscount)

	case model.BxGy:
		d := c.Details.(*model.BxGyDetails)

		count := s.bxgyEligibleCount(cart, d)
		if count == 0 {
			return cart, nil
		}

		for _, g := range d.GetProducts {
			for i := range cart.Items {
				if cart.Items[i].ProductID == g.ProductID {
					freeQty := min(cart.Items[i].Quantity, g.Quantity*count)
					discount := float64(freeQty) * cart.Items[i].Price
					cart.Items[i].TotalDiscount = model.Round(discount)
					cart.TotalDiscount += cart.Items[i].TotalDiscount
				}
			}
		}

		cart.FinalPrice = model.Round(cart.TotalPrice - cart.TotalDiscount)
	}

	return cart, nil
}

func (s *service) estimate(c *model.Coupon, cart *model.Cart) float64 {
	switch c.Type {
	case model.CartWise:
		d := c.Details.(*model.CartWiseDetails)
		if cart.TotalPrice >= d.Threshold {
			return model.Round(cart.TotalPrice * d.Discount / 100)
		}

	case model.ProductWise:
		d := c.Details.(*model.ProductWiseDetails)
		var total float64
		for _, it := range cart.Items {
			if it.ProductID == d.ProductID {
				total += float64(it.Quantity) * it.Price * d.Discount / 100
			}
		}
		return model.Round(total)

	case model.BxGy:
		d := c.Details.(*model.BxGyDetails)
		count := s.bxgyEligibleCount(cart, d)
		if count == 0 {
			return 0
		}
		var total float64
		for _, g := range d.GetProducts {
			for _, it := range cart.Items {
				if it.ProductID == g.ProductID {
					freeQty := min(it.Quantity, g.Quantity*count)
					total += float64(freeQty) * it.Price
				}
			}
		}
		return model.Round(total)
	}

	return 0
}

func (s *service) bxgyEligibleCount(cart *model.Cart, d *model.BxGyDetails) int {
	minGroups := 999999

	for _, req := range d.BuyProducts {
		count := 0
		for _, it := range cart.Items {
			if it.ProductID == req.ProductID {
				count = it.Quantity / req.Quantity
				break
			}
		}
		if count < minGroups {
			minGroups = count
		}
	}

	return min(minGroups, d.RepetitionLimit)
}

func min(a, b int) int {
	if a < b { return a }
	return b
}
