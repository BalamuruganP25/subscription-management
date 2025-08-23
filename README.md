# Subscription Management Service

A simple RESTful API for managing users in a subscription system, built with Go, [chi router](https://github.com/go-chi/chi), and PostgreSQL.

---

## Prerequisites

- [Docker](https://www.docker.com/) (for both Go and PostgreSQL)

---

## Getting Started

1. **Clone the repository:**
   ```sh
   git clone https://github.com/yourusername/subscription-management.git
   cd subscription-management
   ```

2. **Use the Makefile for common tasks:**

   - **🔧 Build Docker containers**
     ```sh
     make build
     ```
   - **🚀 Run the full Docker stack**
     ```sh
     make run
     ```
   - **🧹 Stop and remove containers**
     ```sh
     make stop
     ```
   - **📦 Install and tidy Go dependencies**
     ```sh
     make dep
     ```

   The server will be available at:  
   [http://localhost:8089](http://localhost:8089)

---

## API Endpoints

### 1. Create User 

**POST** `/v1/api/user`

**Request Body:**
```json
{
    "name":"bala",
    "email_id":"bala@gmail.com",
    "phone_number":"9916588437"
}
```

**Response:**
```json
{
    "id":"1",
    "name":"bala",
    "email_id":"bala@gmail.com",
    "phone_number":"99165XXXXX"
}
```

### 2. Get User Details

**GET** `/v1/api/{id}`

**Response:**
```json
{
    "id":"1",
    "name":"bala",
    "email_id":"bala@gmail.com",
    "phone_number":"99165XXXXX",
    "status":"ACTIVE"
}
```

### 3. Delete User Details

**DELETE** `/v1/api/{id}`

**Response:**
```json
{
    "message":"user deleted successfully"
}
```

### 4. Update User Details

**PATCH** `/v1/api/{id}`


**Request Body:**
```json
{
    "phone_number":"9916588437"
}
```

**Response:**
```json
{
    "message":"user deleted successfully"
}
```

### 5. Create Customer

**POST** `/v1/api/customers`

**Request Body:**
```json
{
    "name":"bala",
    "email_id":"bala@gmail.com",
    "phone_number":"9916588437"
}
```

**Response:**
```json
{
    "id":"1",
    "name":"bala",
    "email_id":"bala@gmail.com",
    "phone_number":"99165XXXXX"
}
```

### 6. Create Customer

**POST** `/v1/api/customers`

**Request Body:**
```json
{
    "name":"bala",
    "email_id":"bala@gmail.com",
    "phone_number":"9916588437"
}
```

**Response:**
```json
{
    "id":"1",
    "name":"bala",
    "email_id":"bala@gmail.com",
    "phone_number":"99165XXXXX"
}
```

### 6. Create Subscriptions

**POST** `/v1/api/subscriptions`

**Request Body:**
```json
{
    "customer_id":"bala",
    "price_id":"bala@gmail.com",
    "promo_code":"9916588437"
}
```

**Response:**
```json
{
    "id":"1",
    "name":"bala",
    "email_id":"bala@gmail.com",
    "phone_number":"99165XXXXX",
    "subscription_id":"",
    "subscription_status":""
}
```

### 6. Get Tax details

**GET** `/v1/api/tax/{country}/{state}/{amount}`

**Response:**
```json
{
    "country":"US",
    "state":"NY",
    "tax_rate":10.3,
    "tax_amount":89,
    "amount":1200
}
```

### 6. Create Subscriptions

**POST** `/v1/api/webhooks/stripe`

**Request Body:**
```json
   {
  "id": "",
  "object": "event",
  "type": "customer.subscription.updated",
  "data": {
    "object": {
      "id": "",
      "customer": "",
      "status": "",
      "items": {
        "data": [
          {
            "id": "",
            "price": {
              "id": "price_1RzA35J3MG5r87tLK8obNKk3",
              "unit_amount": 0,
              "currency": "usd"
            }
          }
        ]
      },
      "current_period_end": 1755981731
    }
  }
}

```