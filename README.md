## 📐 Architecture Overview

This template follows **Hexagonal Architecture** (Ports & Adapters). The flow of a request ensures a strict separation of concerns.

### The Request Lifecycle
1.  **Middleware** (Logger/Auth) runs.
2.  Request hits `interfaces/http/router`.
3.  **Handler** parses JSON to DTO.
4.  Handler calls `application/usecase`.
5.  **Usecase** validates business rules.
6.  Usecase calls `domain/repository` method.
7.  **Infrastructure** (implementation) executes SQL.
8.  Usecase returns **Domain Entity**.
9.  **Handler** converts Entity to Response DTO.
10. **Response** sent.

### The Roles
*   **Handler:** Waiter (Passes data).
*   **Middleware:** Bouncer (Checks auth).
*   **DTO:** Translator (JSON <-> Struct).
*   **Application:** Manager (Decides *what* to do).
*   **Domain:** Rulebook (Defines *what* things are).
*   **Infrastructure:** Worker (Does the actual work/SQL).



## 🏗️ Project Structure

```text
.
├── api                    # API Definitions (OpenAPI/Swagger, Proto files)
├── cmd                    # Application Entry Points
│   └── server
│       └── main.go        # The only place where initialization happens (Wiring)
├── configs                # Configuration files
│   └── config.yaml
├── internal               # Private Application Code (The Core)
│   ├── application        # USE CASES (Orchestrators)
│   │   └── service
│   ├── domain             # ENTERPRISE RULES (The Heart)
│   │   ├── entity         # Pure business structs
│   │   ├── repository     # INTERFACES (Ports) - Contracts only!
│   │   └── service        # Domain Services
│   ├── infrastructure     # EXTERNAL TOOLS (Adapters)
│   │   ├── database       # DB connection logic
│   │   └── persistence    # SQL Implementations (Repo Impl)
│   └── interfaces         # DELIVERY MECHANISMS
│       ├── http
│       │   ├── handler    # HTTP Controllers
│       │   ├── middleware
│       │   └── router     # Route definitions (Gin)
│       └── dto            # Data Transfer Objects
├── pkg                    # Public Libraries (Reusable code)
│   ├── errors
│   ├── logger
│   └── utils
├── scripts                # Automation
│   ├── build.sh
│   └── migrate.sh
├── .gitignore
├── go.mod
├── Makefile
└── README.md
```



## 🛡️ SOLID Principles in Action

This template isn't just organized; it's engineered to satisfy the 5 SOLID principles.

### 1. Single Responsibility (SRP)
*   **Rule:** Each folder has one job.
*   **Here:** `domain` holds rules, `interfaces` handles transport, `infrastructure` handles details.
*   **Benefit:** If SQL syntax changes, you only touch `infrastructure`. You never touch `interfaces`.

### 2. Open/Closed (OCP)
*   **Rule:** Software entities should be open for extension but closed for modification.
*   **Here:** You can add a new `user_redis.go` in infrastructure to implement caching without changing the Domain logic or existing `user_postgres.go`.

### 3. Liskov Substitution (LSP)
*   **Rule:** Objects should be replaceable with instances of their subtypes without altering the correctness of the program.
*   **Here:** `MockUserRepository` can replace `PostgresUserRepository` in unit tests seamlessly because they both satisfy the `domain/repository` Interface.

### 4. Interface Segregation (ISP)
*   **Rule:** Clients should not be forced to depend on interfaces they do not use.
*   **Here:** Domain interfaces are small. Instead of one giant `UserRepository`, you can have `UserReader` and `UserWriter` if needed.

### 5. Dependency Inversion (DIP)
*   **Rule:** Depend on abstractions, not concretions.
*   **Here:** `application` layer depends on `domain/repository` (Abstractions), NOT `infrastructure/postgres` (Concretions). The dependency points inward toward the domain.


## ⚙️ Getting Started

### Prerequisites
*   Go 1.21+
*   Supabase Account (or Postgres DB)

### Installation

1.  **Clone the repo**
    ```bash
    git clone https://github.com/cureerel/gotemplate.git
    cd gotemplate
    ```

2.  **Install Dependencies**
    ```bash
    go mod download
    ```

3.  **Configure Environment**
    *   Copy `configs/config.yaml`.
    *   Update the `database.dsn` with your Supabase Connection String (Transaction Mode recommended).

4.  **Run the Server**
    *   Using the Makefile (Recommended):
        ```bash
        make run
        ```
    *   Or manually:
        ```bash
        go run cmd/server/main.go
        ```

The server will start on `0.0.0.0:8080`.



## 🔧 Development

### Makefile Commands
The `Makefile` acts like a task runner (similar to `npm run` in Node.js).

*   `make run`: Runs the application in development mode.
*   `make build`: Compiles the binary to `./build/`.
*   `make clean`: Removes build artifacts.

### Adding a New Feature
1.  **Define Entity:** Add struct to `internal/domain/entity`.
2.  **Define Interface:** Add methods to `internal/domain/repository`.
3.  **Implement Repo:** Write SQL in `internal/infrastructure/persistence`.
4.  **Write Usecase:** Add business logic to `internal/application/service`.
5.  **Create Handler:** Add HTTP logic to `internal/interfaces/http/handler`.
6.  **Register Route:** Add endpoint in `internal/interfaces/http/router`.


## ⚠️ Troubleshooting & Lessons Learned

This template was forged in fire. Here are the specific problems encountered during development and their solutions.

### 1. IPv6 Connection Refused (Supabase)
*   **Problem:** `dial tcp [2406:...]:5432: connect: no route to host`.
*   **Cause:** The local machine or ISP preferred IPv6 DNS resolution, but Supabase's free tier direct connection often drops IPv6 packets.
*   **Solution:** Used the **Transaction Mode** Connection Pooler from Supabase Dashboard (port `6543`). This guarantees a stable IPv4 connection.

### 2. GORM Driver Mismatch
*   **Problem:** `missing method ExecContext` when trying to use `pgx/v5` with GORM.
*   **Cause:** Mixing the raw `pgx` driver with GORM's `postgres` driver caused type incompatibility.
*   **Solution:** Stuck to the standard `gorm.io/driver/postgres` for simplicity and compatibility, using `GODEBUG=netdns=go` to handle DNS resolution issues.

### 3. Makefile "Missing Separator"
*   **Problem:** `Makefile:10: *** missing separator. Stop.`
*   **Cause:** Makefiles require **Tabs** for indentation, not spaces.
*   **Solution:** Ensured the editor used "Indent with Tabs" for the Makefile.

### 4. Import Paths & Module Structure
*   **Problem:** `module github.com/...: git ls-remote ... Repository not found`.
*   **Cause:** Trying to import local folders using full URL paths before the repo existed on GitHub, or misconfigured `go.mod`.
*   **Solution:** Ensured all imports used the module name defined in `go.mod` (e.g., `github.com/cureerel/gotemplate/internal/...`) and ran `go mod tidy`.
