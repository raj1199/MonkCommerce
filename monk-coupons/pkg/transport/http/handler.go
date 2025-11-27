package http

import (
	"context"
	"encoding/json"
	"net/http"

	khttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/you/monk-coupons/pkg/endpoints"
	"github.com/you/monk-coupons/pkg/model"
)

func MakeHTTPHandler(e endpoints.Endpoints) http.Handler {
	r := mux.NewRouter()

	r.Methods("POST").Path("/coupons").Handler(khttp.NewServer(
		e.Create, decodeCoupon, encodeJSON))

	r.Methods("GET").Path("/coupons").Handler(khttp.NewServer(
		e.GetAll, decodeNone, encodeJSON))

	r.Methods("GET").Path("/coupons/{id}").Handler(khttp.NewServer(
		e.GetByID, decodeID, encodeJSON))

	r.Methods("PUT").Path("/coupons/{id}").Handler(khttp.NewServer(
		e.Update, decodeUpdateCoupon, encodeJSON))

	r.Methods("DELETE").Path("/coupons/{id}").Handler(khttp.NewServer(
		e.Delete, decodeID, encodeJSON))

	r.Methods("POST").Path("/applicable-coupons").Handler(khttp.NewServer(
		e.Applicable, decodeCart, encodeJSON))

	r.Methods("POST").Path("/apply-coupon/{id}").Handler(khttp.NewServer(
		e.Apply, decodeApply, encodeJSON))

	return r
}

// --- Decoders ---

func decodeCoupon(_ context.Context, r *http.Request) (interface{}, error) {
	var c model.Coupon
	return &c, json.NewDecoder(r.Body).Decode(&c)
}

func decodeID(_ context.Context, r *http.Request) (interface{}, error) {
	return mux.Vars(r)["id"], nil
}

func decodeUpdateCoupon(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["id"]
	var c model.Coupon
	_ = json.NewDecoder(r.Body).Decode(&c)
	return map[string]interface{}{"id": id, "payload": &c}, nil
}

func decodeCart(_ context.Context, r *http.Request) (interface{}, error) {
	var cart model.Cart
	return &cart, json.NewDecoder(r.Body).Decode(&cart)
}

func decodeApply(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)["id"]
	var cart model.Cart
	_ = json.NewDecoder(r.Body).Decode(&cart)
	return map[string]interface{}{"id": id, "cart": &cart}, nil
}

// --- Encoder ---

func encodeJSON(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func decodeNone(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}
