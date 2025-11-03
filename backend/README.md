# CPRoom Backend Microservices

This repository contains the backend microservices for the CPRoom project.  
The project uses a **hybrid microservice structure** with **REST + gRPC** services and shared proto definitions.

---

## 📂 Project Structure

```
backend/
├── libs/                        # Shared libs
│   ├── log
│   │   └── zap_logger.go        # zap config
│   └── middleware
│   │   └── /...                 # Middleware file for authentication
├── proto/                       # Shared gRPC proto definitions
│   ├── approval.proto
│   ├── booking.proto
│   ├── notification.proto
│   └── user.proto
├── services/
│   ├── auth/                    # Auth Service (REST)
│   │   ├── cmd/
│   │   │   └── server.go        # Fiber server entrypoint
│   │   └── internal/
│   │       ├── auth_handler.go  # REST handlers
│   │       ├── auth_service.go  # Business logic
│   │       └── auth_repo.go     # Repository (DB layer)
│   ├── booking/                 # Booking Service (gRPC)
│   │   └── internal/...
│   ├── notification/            # Notification Service (gRPC)
│   │   └── internal/...
│   ├── user/                    # User Service (gRPC)
│   │   └── internal/...
│   └── approval/                # Staff Approval Service (gRPC)
│       └── internal/...
├── go.mod                       # Top-level Go module
└── README.md
```

---

## ⚙️ Key Points

-   **Hybrid Microservices**

    -   `Auth Service` → REST (Fiber)
    -   Other services (`User`, `Booking`, `Notification`, `Approval`) → gRPC
    -   `API Gateway` (optional) → REST → gRPC translation

-   **Shared Protos**

    -   All proto files live under `backend/proto/`
    -   `go_package` in all proto files:

    ```proto
    option go_package = "github.com/JJnvn/Software-Arch-CPRoom/backend/proto;proto";
    ```

    -   This allows **all services to import shared proto definitions**:

    ```go
    import proto "github.com/JJnvn/Software-Arch-CPRoom/backend/proto"
    ```

-   **Top-level Go module**

    -   `go mod init github.com/JJnvn/Software-Arch-CPRoom/backend`
    -   Ensures all services can import shared proto code consistently.

-   **API Gateway (Kong)**

    -   DB-less Kong gateway lives under `backend/kong/`
    -   Declarative config `kong.yml` wires services, routes, and JWT enforcement
    -   Kong Manager OSS is exposed on `http://localhost:8002` (and `https://localhost:8445`)
    -   Proxy traffic via `https://localhost:8443`; JWTs issued by the auth service are validated at the edge

-   **Auth Service**

    -   Handles `register`, `login`, and `validate` endpoints via REST
    -   Uses JWT for authentication
    -   Repository, service, and handler layers are in `internal/`

-   **User Service**
    -   gRPC service managing profiles, preferences, and booking history
    -   Other services or API Gateway can call it via gRPC

---

## 🚀 Running the Auth Service

```bash
cd backend/services/auth
go run cmd/server.go
```

Endpoints:

-   `POST /register` → Register new user
-   `POST /login` → Login and get JWT
-   `GET /validate` → Validate JWT

---
