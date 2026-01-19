# Butik API Documentation

This document describes the REST API endpoints for the Butik backend service. All endpoints return JSON responses. Use best practices for request validation and error handling as shown below.

---

## Authentication

### 1. Login

- **POST** `/login`
- **Description:** Authenticate admin user and get JWT tokens.
- **Request Body:**
  | Field | Type | Required | Validation |
  |----------|--------|----------|---------------------------|
  | username | string | Yes | min:3, max:50 |
  | password | string | Yes | min:6, max:100 |
- **Example:**

```json
{
  "username": "admin",
  "password": "yourpassword"
}
```

- **Response:**

```json
{
  "message": "Login successful",
  "access_token": "...",
  "refresh_token": "..."
}
```

### 2. Refresh Token

- **POST** `/refresh-token`
- **Description:** Get a new access token using a refresh token.
- **Request Body:**
  | Field | Type | Required | Validation |
  |--------------|--------|----------|------------|
  | refresh_token| string | Yes | required |
- **Example:**

```json
{
  "refresh_token": "..."
}
```

- **Response:**

```json
{
  "access_token": "..."
}
```

---

## Product

### 1. List Products

- **GET** `/products?page=1&limit=10`
- **Description:** Get paginated list of products.
- **Query Params:**
  - `page` (int, optional, default: 1)
  - `limit` (int, optional, default: 10)
- **Response:**

```json
{
  "data": [ ...products ],
  "total": 100,
  "page": 1,
  "limit": 10
}
```

### 2. Get Product by ID

- **GET** `/products/{id}`
- **Description:** Get product details by ID.
- **Response:**

```json
{
  "id": 1,
  "name": "Product Name",
  "description": "...",
  "price": 10000,
  "stock": 10,
  "category": { ... },
  "image_url": "...",
  "created_at": "..."
}
```

### 3. Create Product

- **POST** `/products` (Protected, JWT)
- **Description:** Create a new product. Use `multipart/form-data` for image upload.
- **Form Fields:**
  | Field | Type | Required | Validation |
  |-------------|---------|----------|---------------------------|
  | name | string | Yes | min:2, max:200 |
  | description | string | Yes | min:10, max:2000 |
  | price | float | Yes | gt:0, lte:999999999 |
  | stock | int | Yes | gte:0, lte:99999 |
  | category_id | uint | Yes | gt:0 |
  | image | file | Yes | image file |
- **Response:**

```json
{
  "message": "Product created successfully",
  "product": { ... }
}
```

### 4. Update Product

- **PUT** `/products/{id}` (Protected, JWT)
- **Description:** Update product info. Use `multipart/form-data` for image upload (optional).
- **Form Fields:** (same as Create Product)
- **Response:**

```json
{
  "message": "Product updated successfully",
  "product": { ... }
}
```

### 5. Delete Product

- **DELETE** `/products/{id}` (Protected, JWT)
- **Description:** Delete a product by ID.
- **Response:**

```json
{
  "message": "Product deleted successfully"
}
```

---

## Category

### 1. List Categories

- **GET** `/categories?page=1&limit=10`
- **Description:** Get paginated list of categories.
- **Query Params:**
  - `page` (int, optional, default: 1)
  - `limit` (int, optional, default: 10)
- **Response:**

```json
{
  "data": [ ...categories ],
  "page": 1,
  "limit": 10,
  "total": 20
}
```

### 2. Get Category by ID

- **GET** `/categories/{id}`
- **Description:** Get category details by ID.
- **Response:**

```json
{
  "id": 1,
  "name": "Category Name",
  "created_at": "..."
}
```

### 3. Create Category

- **POST** `/categories` (Protected, JWT)
- **Description:** Create a new category.
- **Request Body:**
  | Field | Type | Required | Validation |
  |-------|--------|----------|--------------------|
  | name | string | Yes | min:2, max:100 |
- **Response:**

```json
{
  "message": "Category created successfully",
  "category": { ... }
}
```

### 4. Update Category

- **PUT** `/categories/{id}` (Protected, JWT)
- **Description:** Update category name.
- **Request Body:** (same as Create Category)
- **Response:**

```json
{
  "message": "Category updated successfully",
  "category": { ... }
}
```

### 5. Delete Category

- **DELETE** `/categories/{id}` (Protected, JWT)
- **Description:** Delete a category by ID.
- **Response:**

```json
{
  "message": "Category deleted successfully"
}
```

---

## Order

### 1. Create Order

- **POST** `/orders`
- **Description:** Create a new order. Use `multipart/form-data` for proof of payment upload.
- **Form Fields:**
  | Field | Type | Required | Validation |
  |-----------------|---------|----------|---------------------------|
  | customer_name | string | Yes | min:2, max:100 |
  | whatsapp | string | Yes | min:10, max:15 |
  | map_address | string | No | max:500 |
  | latitude | float | No | gte:-90, lte:90 |
  | longitude | float | No | gte:-180, lte:180 |
  | address_note | string | No | max:500 |
  | items | JSON | Yes | array of order items |
  | proof_of_payment| file | Yes | image file |
- **Order Item Format:**

```json
[{ "product_id": 1, "quantity": 2 }]
```

- **Response:**

```json
{
  "message": "Order created successfully",
  "order": { ... }
}
```

### 2. Get Order by ID

- **GET** `/orders/{id}`
- **Description:** Get order details by order ID.
- **Response:**

```json
{
  "id": "...",
  "customer_name": "...",
  "status": "pending",
  ...
}
```

### 3. List Orders (Admin)

- **GET** `/orders?page=1&limit=10` (Protected, JWT)
- **Description:** Get paginated list of all orders.
- **Query Params:**
  - `page` (int, optional, default: 1)
  - `limit` (int, optional, default: 10)
- **Response:**

```json
{
  "data": [ ...orders ],
  "page": 1,
  "limit": 10,
  "total": 50
}
```

### 4. Update Order Status (Admin)

- **PUT** `/orders/{id}/status` (Protected, JWT)
- **Description:** Update order status (success/reject/pending).
- **Request Body:**
  | Field | Type | Required | Validation |
  |--------|-------------|----------|-----------------------------|
  | status | string enum | Yes | one of: pending, success, reject |
- **Example:**

```json
{
  "status": "success"
}
```

- **Response:**

```json
{
  "message": "Order status updated successfully",
  "order": { ... }
}
```

### 5. Delete Order (Admin)

- **DELETE** `/orders/{id}` (Protected, JWT)
- **Description:** Delete an order by ID.
- **Response:**

```json
{
  "message": "Order deleted successfully"
}
```

---

## Error Response Format

All error responses use this format:

```json
{
  "error": "Error message here"
}
```

---

## Notes

- All protected endpoints require JWT in the `Authorization: Bearer <token>` header.
- Use proper validation for all input fields as described above.
- File uploads use `multipart/form-data`.
- Pagination is supported for list endpoints via `page` and `limit` query params.
