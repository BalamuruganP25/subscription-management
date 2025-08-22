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

| Method | Endpoint           | Description         |
|--------|--------------------|---------------------|
| POST   | `/v1/user`         | Create a new user   |
| GET    | `/v1/user/{id}`    | Get user by ID      |
| PATCH  | `/v1/user/{id}`    | Update user by ID   |
| DELETE | `/v1/user/{id}`    | Delete