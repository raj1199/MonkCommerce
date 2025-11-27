package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"
)

type CouponType string

const (
	CartWise    CouponType = "cart-wise"
	ProductWise CouponType = "product-wise"
	BxGy        CouponType = "bxgy"
)

type Coupon struct {
	ID              string      `json:"id" bson:"_id,omitempty"`
	Type            CouponType  `json:"type" bson:"type"`
	Details         interface{} `json:"details" bson:"details"`
	CreatedAt       time.Time   `json:"created_at" bson:"created_at"`
	ExpiresAt       *time.Time  `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
	RepetitionLimit int         `json:"repetition_limit,omitempty" bson:"repetition_limit,omitempty"`
}

type CartWiseDetails struct {
	Threshold float64 `json:"threshold" bson:"threshold"`
	Discount  float64 `json:"discount" bson:"discount"`
}

type ProductWiseDetails struct {
	ProductID int     `json:"product_id" bson:"product_id"`
	Discount  float64 `json:"discount" bson:"discount"`
}

type ProductQuantity struct {
	ProductID int `json:"product_id" bson:"product_id"`
	Quantity  int `json:"quantity" bson:"quantity"`
}

type BxGyDetails struct {
	BuyProducts     []ProductQuantity `json:"buy_products" bson:"buy_products"`
	GetProducts     []ProductQuantity `json:"get_products" bson:"get_products"`
	RepetitionLimit int               `json:"repetition_limit" bson:"repetition_limit"`
}

type CartItem struct {
	ProductID     int     `json:"product_id"`
	Quantity      int     `json:"quantity"`
	Price         float64 `json:"price"`
	TotalDiscount float64 `json:"total_discount,omitempty"`
}

type Cart struct {
	Items         []CartItem `json:"items"`
	TotalPrice    float64    `json:"total_price,omitempty"`
	TotalDiscount float64    `json:"total_discount,omitempty"`
	FinalPrice    float64    `json:"final_price,omitempty"`
}

type ApplicableResponse struct {
	CouponID string  `json:"coupon_id"`
	Type     string  `json:"type"`
	Discount float64 `json:"discount"`
}

func (c *Cart) CalcTotal() {
	var total float64
	for _, it := range c.Items {
		total += it.Price * float64(it.Quantity)
	}
	c.TotalPrice = Round(total)
	c.FinalPrice = Round(total - c.TotalDiscount)
}

func Round(n float64) float64 {
	return math.Round(n*100) / 100
}

// Validation Rules

var (
	ErrInvalidDiscount    = errors.New("discount must be between 0 and 100")
	ErrInvalidThreshold   = errors.New("threshold cannot be negative")
	ErrInvalidProductQty  = errors.New("product quantity must be >= 1")
	ErrInvalidCouponType  = errors.New("invalid coupon type")
	ErrInvalidBxGyPayload = errors.New("invalid bxgy payload")
)

func ValidateCoupon(c *Coupon) error {
	switch c.Type {
	case CartWise:
		d1, _ := json.Marshal(c.Details)
	    var d CartWiseDetails
		err := json.Unmarshal(d1, &d)
		if err != nil {
			fmt.Println("Error in cart-wise - ",ErrInvalidCouponType)
			return ErrInvalidCouponType
		}
		if d.Threshold < 0 {
			fmt.Println("Error in cart-wise ",ErrInvalidThreshold)
			return ErrInvalidThreshold
		}
		if d.Discount < 0 || d.Discount > 100 {
			fmt.Println("Error in cast-wise ",ErrInvalidDiscount)
			return ErrInvalidDiscount
		}

	case ProductWise:
		d1, _ := json.Marshal(c.Details)
	    var d ProductWiseDetails
		err := json.Unmarshal(d1, &d)
		if err != nil {
			fmt.Println("Error in product-wise - ",ErrInvalidCouponType)
			return ErrInvalidCouponType
		}
		if d.ProductID <= 0 {
			fmt.Println("Error in product-wise ",ErrInvalidProductQty)
			return ErrInvalidProductQty
		}
		if d.Discount < 0 || d.Discount > 100 {
			fmt.Println("Error in product-wise ",ErrInvalidDiscount)
			return ErrInvalidDiscount
		}

	case BxGy:
		d1, _ := json.Marshal(c.Details)
	    var d BxGyDetails
		err := json.Unmarshal(d1, &d)
		if err != nil {
			fmt.Println("Error in BxGy - ",ErrInvalidCouponType)
			return ErrInvalidCouponType
		}
		if d.RepetitionLimit < 0 {
			fmt.Println("Error in BxGy - ",ErrInvalidBxGyPayload)
			return ErrInvalidBxGyPayload
		}
		for _, b := range d.BuyProducts {
			if b.ProductID <= 0 || b.Quantity <= 0 {
				fmt.Println("Error in BxGy - ",ErrInvalidProductQty)
				return ErrInvalidProductQty
			}
		}
		for _, g := range d.GetProducts {
			if g.ProductID <= 0 || g.Quantity <= 0 {
				fmt.Println("Error in BxGy - ",ErrInvalidProductQty)
				return ErrInvalidProductQty
			}
		}

	default:
		fmt.Println("Invalid coupon-type")
		return ErrInvalidCouponType
	}

	return nil
}
