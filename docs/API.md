# Polygame API Documentation

This document provides a detailed specification for the Polygame API.

**Base URL**: `/api/v1`

## Authentication

All authenticated endpoints require a `Bearer` token in the `Authorization` header.

`Authorization: Bearer <your_jwt_token>`

---

## 1. Auth Endpoints

### 1.1 Register

- **Endpoint**: `POST /auth/register`
- **Description**: Creates a new user account.
- **Request Body**:

```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

- **Success Response (201 Created)**:

```json
{
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "virtual_balance": 10000,
    "is_admin": false,
    "created_at": "2026-01-07T00:00:00Z"
  },
  "token": "jwt_token_string"
}
```

- **Error Responses**:
  - `400 Bad Request`: Invalid input, or username/email already exists.

### 1.2 Login

- **Endpoint**: `POST /auth/login`
- **Description**: Authenticates a user and returns a JWT token.
- **Request Body**:

```json
{
  "username": "testuser",
  "password": "password123"
}
```

- **Success Response (200 OK)**:

```json
{
  "user": { ... },
  "token": "jwt_token_string"
}
```

- **Error Responses**:
  - `401 Unauthorized`: Invalid credentials.

---

## 2. User Endpoints

*(Authentication required)*

### 2.1 Get Profile

- **Endpoint**: `GET /user/profile`
- **Description**: Retrieves the profile of the currently authenticated user.

### 2.2 Update Profile

- **Endpoint**: `PUT /user/profile`
- **Description**: Updates the user's profile (e.g., avatar).
- **Request Body**:

```json
{
  "avatar": "https://example.com/avatar.png"
}
```

### 2.3 Get Balance

- **Endpoint**: `GET /user/balance`
- **Description**: Retrieves the user's current virtual balance.

---

## 3. Market Endpoints

*(Authentication required)*

### 3.1 List Markets

- **Endpoint**: `GET /markets`
- **Description**: Retrieves a paginated list of markets.
- **Query Parameters**:
  - `category` (string, optional): Filter by category (e.g., `sports`).
  - `status` (string, optional): Filter by status (`active`, `closed`, `resolved`).
  - `page` (int, optional): Page number (default: 1).
  - `page_size` (int, optional): Items per page (default: 20).

### 3.2 Get Market Details

- **Endpoint**: `GET /markets/:id`
- **Description**: Retrieves details for a specific market, including its outcomes.

### 3.3 Get Trending Markets

- **Endpoint**: `GET /markets/trending`
- **Description**: Retrieves a list of trending markets, sorted by trading volume.

### 3.4 Search Markets

- **Endpoint**: `GET /markets/search`
- **Description**: Searches for markets by keyword.
- **Query Parameters**:
  - `q` (string, required): The search keyword.

---

## 4. Trading Endpoints

*(Authentication required)*

### 4.1 Place Order

- **Endpoint**: `POST /trading/orders`
- **Description**: Places a buy or sell order for an outcome share.
- **Request Body**:

```json
{
  "market_id": 1,
  "outcome_id": 1,
  "order_type": "buy",
  "shares": 10,
  "price": 0.6
}
```

### 4.2 Get User Orders

- **Endpoint**: `GET /trading/orders`
- **Description**: Retrieves a paginated list of the user's orders.

### 4.3 Get User Positions

- **Endpoint**: `GET /trading/positions`
- **Description**: Retrieves the user's current positions across all markets.

### 4.4 Cancel Order

- **Endpoint**: `DELETE /trading/orders/:id`
- **Description**: Cancels a pending order.

---

## 5. Admin Endpoints

*(Admin authentication required)*

### 5.1 List Users

- **Endpoint**: `GET /admin/users`
- **Description**: Retrieves a paginated list of all users.

### 5.2 Create Market

- **Endpoint**: `POST /admin/markets`
- **Description**: Creates a new market.
- **Request Body**:

```json
{
  "title": "New Market Title",
  "description": "Market description.",
  "category": "sports",
  "image_url": "https://example.com/image.jpg",
  "outcomes": ["Outcome A", "Outcome B"]
}
```

### 5.3 Update Market

- **Endpoint**: `PUT /admin/markets/:id`
- **Description**: Updates an existing market's details.

### 5.4 Resolve Market

- **Endpoint**: `POST /admin/markets/:id/resolve`
- **Description**: Resolves a market and settles all positions.
- **Request Body**:

```json
{
  "winning_outcome_id": 1
}
```
