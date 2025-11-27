Absolutely — below is your final **Postman-style curl collection + README.md** matching your implemented logic and unimplemented future cases.

You can **copy-paste the README directly into your GitHub repo**.

---

---

# 📘 README.md — Monk Commerce Coupon Service

### 🏗 Built With

* **Golang 1.22**
* **Go-Kit Microservice Architecture**
* **MongoDB**
* **REST API**
* **Typed Validation + Rule Engine for Coupons**

---

## 📌 Project Purpose

This backend system manages coupon creation, evaluation, and application logic with support for:

| Coupon Type           | Supported? | Notes                                                |
| --------------------- | ---------- | ---------------------------------------------------- |
| Cart-wise Discount    | ✅ Yes      | Based on minimum cart value                          |
| Product-wise Discount | ✅ Yes      | Applies only to target product(s)                    |
| BxGy (Buy X Get Y)    | ✅ Yes      | Supports repeated eligibility with repetition limits |

---

---

## ✔️ Implemented Test Cases

| Case Category            | Test Case                                      | Behavior                                   |
| ------------------------ | ---------------------------------------------- | ------------------------------------------ |
| **Cart-wise Coupons**    | Apply discount when cart total meets threshold | Discount = % of cart total                 |
|                          | Reject coupon when below threshold             | No discount applied                        |
|                          | Handle rounding precision                      | Uses `math.Round()` for 2-decimal accuracy |
| **Product-wise Coupons** | Apply % discount to specific product only      | Per-item discount multiplied by quantity   |
|                          | If item not present in cart → ignore           | Coupon still valid but discount = 0        |
| **BxGy Coupons**         | Buy X items get Y free                         | Free quantity depends on repetition limit  |
|                          | Supports multiple repetition cycles            | e.g., Buy 2 get 1 free → 6 items ⇒ 3 free  |
|                          | Handles capped free quantities                 | Free cannot exceed available quantity      |
| **Global Rules**         | Expired coupons are rejected                   | Based on `ExpiresAt` timestamp             |
|                          | Applies only one coupon at a time              | No stacking                                |
|                          | Returns sorted applicable coupons              | API: `/applicable-coupons`                 |

---

---

## 🚧 Not Implemented (Future Scope) — and How We Would Implement

| Unimplemented Case                                  | Why Not Implemented                             | Future Implementation Strategy                                          |
| --------------------------------------------------- | ----------------------------------------------- | ----------------------------------------------------------------------- |
| Coupon stacking (multiple coupons applied together) | Requires conflict rules, precedence, exclusions | Add `priority` + `exclusive` flags in schema and modify rule evaluation |
| Category-based coupons                              | Requires product catalog or metadata            | Add `category` field and join lookup with catalog API                   |
| User-based rules (first-order, max usage per user)  | Needs authentication layer                      | Add user context → store usage history in MongoDB                       |
| Max discount cap (e.g., "50% up to ₹300")           | Business logic extension                        | Add optional field `max_discount`, modify calculation step              |
| Scheduled activation windows (happy hours)          | Requires cron-based activation                  | Add `start_at` + `end_at` fields and scheduled validation               |
| Marketplace vendor restriction                      | Requires integration with merchant catalog      | Add vendorId metadata in products and match with coupon filters         |

---

---

## 🚀 Running the Project

### 1️⃣ Setup Environment File

Create `.env`:

```
MONGO_URI=mongodb://localhost:27017
DB_NAME=monk_coupons
COLLECTION_NAME=coupons
PORT=8080
```

### 2️⃣ Start Server

```sh
go mod tidy
go run ./cmd/server
```


---

## 🧰 API Endpoints — CURL Collection

---

### 🟩 Create Cart-Wise Coupon

```sh
curl -X POST http://localhost:8080/coupons \
-H "Content-Type: application/json" \
-d '{
  "type": "cart-wise",
  "details": { "threshold": 1000, "discount": 10 }
}'
```

---

### 🟧 Create Product-Wise Coupon

```sh
curl -X POST http://localhost:8080/coupons \
-H "Content-Type: application/json" \
-d '{
  "type": "product-wise",
  "details": { "product_id": 101, "discount": 20 }
}'
```

---

### 🟥 Create BxGy Coupon

```sh
curl -X POST http://localhost:8080/coupons \
-H "Content-Type: application/json" \
-d '{
  "type": "bxgy",
  "details": {
    "buy_products":[{"product_id":101,"quantity":2}],
    "get_products":[{"product_id":102,"quantity":1}],
    "repetition_limit":5
  }
}'
```

---

### 📌 Get All Coupons

```sh
curl http://localhost:8080/coupons
```

---

### 📌 Get Single Coupon

```sh
curl http://localhost:8080/coupons/<id>
```

---

### 📌 Delete Coupon

```sh
curl -X DELETE http://localhost:8080/coupons/<id>
```

---

### 🧮 Get Applicable Coupons for a Cart

```sh
curl -X POST http://localhost:8080/applicable-coupons \
-H "Content-Type: application/json" \
-d '{
  "items":[
    {"product_id":101,"quantity":2,"price":500},
    {"product_id":102,"quantity":1,"price":300}
  ]
}'
```

---

### 🏷 Apply Coupon

```sh
curl -X POST http://localhost:8080/apply-coupon/<coupon-id> \
-H "Content-Type: application/json" \
-d '{
  "items":[
    {"product_id":101,"quantity":2,"price":500},
    {"product_id":102,"quantity":1,"price":300}
  ]
}'
```


## 🧱 Assumptions

* Prices and cart data are trusted inputs.
* Only one coupon can be applied per transaction.
* Calculations require 2-decimal precision.
* Service does not validate product catalog externally.

