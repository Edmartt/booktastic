## Booktastic

**Booktastic** is a personal project for tracking books across different states — acquired, reading, read, and wishlist.
This repository holds the system design documentation: architecture diagrams, sequence diagrams, and design decisions for each microservice. Each microservice has its own README with technical details, API documentation, and setup instructions.

The following diagram shows the overall system architecture:


![](https://github.com/Edmartt/booktastic/blob/dev/assets/system%20architecture-2026.jpg)

## Implementation

The following diagram reflects the concrete tech stack used:
**Go** · **Traefik** · **PostgreSQL** · **Auth0** · **Kafka** · **MongoDB**

![](https://github.com/Edmartt/booktastic/blob/dev/assets/architecture.jpg)


## Running the project

### Prerequisites
- Docker
- Docker Compose
- An Auth0 account with a Regular Web Application configured

### Setup
Create a `.env` file based on `.env.example` and fill in your credentials.

### Start
```bash
docker compose up --build
```

### Services
| Service | URL |
|---------|-----|
| Books API | http://books.localhost |
| Auth API | http://auth.localhost |
| Traefik Dashboard | http://localhost:8080 |
| Swagger - Auth | http://auth.localhost/swagger/index.html |
| Swagger - Books | http://books.localhost/swagger/index.html |

## Authentication Flow

### 1. Register
```bash
curl -X POST http://auth.localhost/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"your@email.com","password":"Password1*","password_confirmation":"Password1*"}'
```

### 2. Login
```bash
curl -X POST http://auth.localhost/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"your@email.com","password":"Password1*"}'
```

### 3. Use the token
```bash
curl -X POST http://books.localhost/api/v1/books \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"isbn":"978-3-16-148410-0","title":"Clean Code","pages":"464","current_page":"0","author":"Robert C. Martin","year":"2008","status":"reading"}'
```
