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

### Installation
```bash
git clone https://github.com/Edmartt/booktastic.git
```
```bash
cd booktastic
```

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

## Books API

### Create a book
```bash
curl -X POST http://books.localhost/api/v1/books \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"isbn":"978-3-16-148410-0","title":"Clean Code","pages":"464","current_page":"0","author":"Robert C. Martin","year":"2008","status":"reading"}'
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Get a book
```bash
curl -X GET http://books.localhost/api/v1/books/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <your_token>"
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "isbn": "978-3-16-148410-0",
  "title": "Clean Code",
  "pages": "464",
  "current_page": "0",
  "author": "Robert C. Martin",
  "year": "2008",
  "status": "reading"
}
```

### Book not found
```json
{
  "message": "book not found"
}
```

### Update a book
```bash
curl -X PUT http://books.localhost/api/v1/books/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"current_page":"100","status":"reading"}'
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "isbn": "978-3-16-148410-0",
  "title": "Clean Code",
  "pages": "464",
  "current_page": "100",
  "author": "Robert C. Martin",
  "year": "2008",
  "status": "reading"
}
```


## Design Decisions

### Authentication

Authentication is handled by Auth0, a third-party identity provider. Rather than implementing credential storage and management from scratch, the decision was to delegate this responsibility to a system specifically built and maintained for that purpose.

This means no passwords are stored in the application's database. Auth0 handles user creation, login, and email confirmation, while the Auth service acts as the integration layer — exposing the `/login` and `/signup` endpoints and validating JWT tokens on behalf of Traefik.

Vendor lock-in was a consideration from the start. The Auth service is designed around interfaces and a factory pattern, keeping the provider-specific implementation behind an abstraction layer. Switching to a different identity provider — such as AWS Cognito — would require a new implementation of that interface, not changes across the codebase.

This approach reduces security risk and maintenance overhead, while still keeping control over what happens after authentication within the system.

### API Gateway

Traefik was chosen as the API gateway, acting as the single entry point for all incoming requests. One of its core responsibilities in this system is JWT validation through a forward auth middleware — every request is verified before it reaches any service, keeping authentication logic centralized and out of individual microservices.

While Traefik is primarily a reverse proxy rather than a fully featured API gateway, it fits the current needs of the project and was already familiar territory. As the system grows, replacing or complementing it with a more complete solution remains an option.

### Book Service & Storage

The book service is the core of the system, responsible for managing the books a user has acquired, is reading, has read, or wants to acquire. It is intentionally kept focused on that single responsibility — no authentication logic, no notification concerns.

The service is structured following hexagonal architecture, keeping the domain logic isolated from external concerns such as the database, the HTTP layer, or any framework. This makes the core logic easier to test, reason about, and maintain — and allows infrastructure details like the database engine to change without touching business logic.

PostgreSQL is used as the primary database in production. For local development, SQLite is used instead, keeping the setup lightweight without external dependencies. The switch between engines is handled through environment variables, with no changes to the application code — a direct benefit of the ports and adapters model that hexagonal architecture enforces.

### Messaging & Event-Driven Notifications

Notifications are handled asynchronously through an event-driven approach using Kafka as the message broker. When a relevant action occurs in the book service — such as saving a book or changing its state — an event is published to a Kafka topic. A dedicated notification worker consumes those events independently and decides how to notify the user.

This design keeps the book service focused on its own responsibility. It has no knowledge of how or whether a notification is sent — it simply publishes an event and moves on. The response is returned to the client immediately, without waiting for the notification to be processed.

The notification worker supports multiple channels — email, Telegram, and SMS — each implemented behind a common interface, following the same principle applied to authentication. Adding or swapping a channel does not affect the rest of the system.

MongoDB is used to persist the notification history, logging the result of each notification sent. This provides visibility into what was delivered, when, and through which channel — useful for debugging and auditing without relying on Kafka's retention window alone.

### Development Notes

This project is actively under development. The messaging and event-driven notifications layer represents one of its main challenges — not only because it involves integrating multiple third-party services, but because it pushes into territory that is still being explored.

Previous experience with workers has been in simpler contexts: a Python worker handling arithmetic operations over geospatial coordinates using Redis Streams, with no third-party service integrations involved. Kafka entered the picture professionally about a year ago in a professional project, and this system is an opportunity to go deeper — designing the full flow from event production to consumption and notification delivery from scratch.


### Pending Decisions


Data model for the book service (ER diagram to be added once finalized)

Active notification channels and provider configuration
