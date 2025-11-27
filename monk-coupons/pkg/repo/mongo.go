package repo

import (
	"context"
	"errors"
	"time"

	"github.com/you/monk-coupons/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo "go.mongodb.org/mongo-driver/mongo"
)

var ErrNotFound = errors.New("coupon not found")

type MongoRepo struct {
	col *mongo.Collection
}

func NewMongoRepo(client *mongo.Client, db, col string) *MongoRepo {
	return &MongoRepo{col: client.Database(db).Collection(col)}
}

func (m *MongoRepo) Create(ctx context.Context, c *model.Coupon) (string, error) {
	c.CreatedAt = time.Now().UTC()

	res, err := m.col.InsertOne(ctx, c)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (m *MongoRepo) GetAll(ctx context.Context) ([]*model.Coupon, error) {
	cur, err := m.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []*model.Coupon
	for cur.Next(ctx) {
		raw := bson.M{}
		_ = cur.Decode(&raw)
		c, _ := hydrateCoupon(raw)
		out = append(out, c)
	}
	return out, nil
}

func (m *MongoRepo) GetByID(ctx context.Context, id string) (*model.Coupon, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": oid}
	raw := bson.M{}
	err := m.col.FindOne(ctx, filter).Decode(&raw)
	if err != nil {
		return nil, ErrNotFound
	}

	return hydrateCoupon(raw)
}

func (m *MongoRepo) Update(ctx context.Context, id string, c *model.Coupon) error {
	oid, _ := primitive.ObjectIDFromHex(id)

	_, err := m.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": c})
	return err
}

func (m *MongoRepo) Delete(ctx context.Context, id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)

	_, err := m.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func hydrateCoupon(raw bson.M) (*model.Coupon, error) {
	c := &model.Coupon{
		ID:              raw["_id"].(primitive.ObjectID).Hex(),
		Type:            model.CouponType(raw["type"].(string)),
		CreatedAt:       raw["created_at"].(primitive.DateTime).Time(),
	}

	if exp, ok := raw["expires_at"]; ok && exp != nil {
		t := exp.(primitive.DateTime).Time()
		c.ExpiresAt = &t
	}

	if exp, ok := raw["repetition_limit"]; ok && exp != nil {
		t := exp.(primitive.DateTime).Time()
		c.ExpiresAt = &t
	}

	detail := raw["details"].(bson.M)
	switch c.Type {
	case model.CartWise:
		var d model.CartWiseDetails
		bs, _ := bson.Marshal(detail)
		_ = bson.Unmarshal(bs, &d)
		c.Details = &d

	case model.ProductWise:
		var d model.ProductWiseDetails
		bs, _ := bson.Marshal(detail)
		_ = bson.Unmarshal(bs, &d)
		c.Details = &d

	case model.BxGy:
		var d model.BxGyDetails
		bs, _ := bson.Marshal(detail)
		_ = bson.Unmarshal(bs, &d)
		c.Details = &d
	}

	return c, nil
}
