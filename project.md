# Cureerel Project Documentation

This document provides a comprehensive guide to the project's architecture, organization, and the workflow for adding or modifying features.

---

## 🏗️ Architecture Overview

The project follows a **Clean Architecture (Layered)** pattern, ensuring a strong separation of concerns.

1.  **Domain Layer (`internal/domain`)**:
    *   **Entities**: Core business objects and logic (e.g., `User`, `Blog`).
    *   **Repositories**: Interfaces defines how entities are stored/retrieved (abstraction).
2.  **Application Layer (`internal/application/service`)**:
    *   **Services**: Implements "Use Cases". Orchestrates entities and calls repository interfaces to fulfill business requirements.
3.  **Infrastructure Layer (`internal/infrastructure`)**:
    *   **Postgres/GORM**: Concrete implementations of Repository interfaces.
    *   **Models**: Database-specific structs (GORM models) and conversion logic to/from Domain Entities.
    *   **External Clients**: Cloudinary, Stripe, Razorpay, Resend.
4.  **Interface Layer (`internal/interfaces`)**:
    *   **HTTP Handlers**: Extracts request data, calls Services, and returns responses.
    *   **DTO (Data Transfer Objects)**: Structs used for request binding and JSON responses.
    *   **Router**: Modular registration of URL paths and middleware.

---

## 🛠️ Feature Addition Workflow

Follow these steps to add a new feature (e.g., "Gift Cards").

### Step 1: Define Domain Entity & Repository
Create the entity in `internal/domain/entity/` and the interface in `internal/domain/repository/`.
*   **Entity**: `type GiftCard struct { ... }`
*   **Repository**: `type GiftCardRepository interface { ... }`

### Step 2: Implement Infrastructure (Postgres)
1.  **Model**: Create a GORM model in `internal/infrastructure/postgres/models/gift_card.go`. Implement `ToDomain()` and `GiftCardFromDomain()` methods.
2.  **Repository Implementation**: Create `internal/infrastructure/postgres/repositories/gift_card_repository.go`. This implements the domain interface using GORM.

### Step 3: Database Migration
Cureerel uses **Atlas** for migrations.
1.  Create a migration file:
    ```bash
    make migrate-create NAME=add_gift_cards
    ```
2.  Edit the generated `.sql` file in `migrations/` to add the `CREATE TABLE` statement.
3.  Calculate the directory hash:
    ```bash
    atlas migrate hash --dir "file://migrations"
    ```
4.  Apply the migration:
    ```bash
    make migrate-up
    ```

### Step 4: Create Application Service
Create `internal/application/service/gift_card_service.go`.
*   Inject the `GiftCardRepository` interface into the service struct.
*   Implement business logic (e.g., `IssueGiftCard`).

### Step 5: Implement HTTP Handler & DTOs
1.  **DTO**: Define request/response structs in `internal/interfaces/dto/`.
2.  **Handler**: Create `internal/interfaces/http/handler/gift_card_handler.go`.
    *   Inject the `GiftCardService` into the handler.
    *   Use `c.ShouldBindJSON` with your DTOs.
    *   Use `helper.go` methods for standardized error/success responses.

### Step 6: Modular Routing
1.  Create `internal/interfaces/http/router/gift_card_routes.go`.
    ```go
    func registerGiftCardRoutes(rg *gin.RouterGroup, d *Deps) {
        gift := rg.Group("/gift-cards")
        gift.Use(middleware.AuthMiddleware(d.AuthService))
        {
            gift.POST("", d.GiftCardHandler.Issue)
        }
    }
    ```
2.  Register the function in `internal/interfaces/http/router/mount_routes.go`.

### Step 7: Wiring in `main.go`
1.  Update `internal/interfaces/http/router/deps.go` to include the new handler.
2.  Update `internal/interfaces/http/router/router.go` to accept the handler.
3.  In `cmd/server/main.go`:
    *   Initialize the **Repository**.
    *   Initialize the **Service**.
    *   Initialize the **Handler**.
    *   Pass the handler to `router.SetupRouter`.

---

## 🛠️ Development Common Commands

*   **Run Server**: `make run` (automatically loads `.env`).
*   **Create Migration**: `make migrate-create NAME=title`
*   **Apply Migration**: `make migrate-up`
*   **Force Migration**: `make migrate-force VERSION=timestamp` (to fix out-of-sync schemas).

## 🧩 Connection Summary (The "Why")

*   **Handler ↔ DTO**: Handlers use DTOs to validate input. This keeps "scraping" logic out of services.
*   **Service ↔ Entity**: Services only speak in "Entities". They don't know about JSON or Databases.
*   **Repository ↔ Model**: Repositories translate safe "Domain Entities" into "DB Models" for persistence.

By following this flow, you can modify any part of the system (e.g., changing Stripe to another provider) by only touching the **Infrastructure** and **Wiring** layers, while leaving your **Business Logic** (Service/Domain) untouched.
