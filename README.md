# Project Title
CProom - Meeting Room Booking Application

## Getting Started
Frontend project using Next.js.
Backend project using Golang.

### Prerequisites
Install 
[Docker](https://www.docker.com/get-started) to run the whole project using docker-compose.
[Bun](https://bun.sh/) for frontend and 
[Go](https://golang.org/dl/) for backend.

In windows, you may need to enable WSL2 and install a Linux distribution from the Microsoft Store.
because Next.js hot-reloading may not work properly in Windows file system.

### Installation
Install the dependencies for both frontend and backend.
for frontend:
```bash
cd frontend/cproom
bun install
```
for backend:
No dependencies to install.


## Usage

To run Whole project:
```bash
docker compose up -d --build
```
To run Frontend only:
```bash
cd frontend/cproom
bun dev
```
To run Backend only:
```bash
cd backend
air
```

-   API Gateway (Kong) proxies the backend on `https://localhost:8443` with a self-signed certificate located in `backend/kong/certs`. Kong Manager OSS is available on `http://localhost:8002`.
-   The web UI is now accessible over HTTPS at `https://localhost:3443` (same certificate).

### API Specifications

-   `backend/services/auth/openapi.yaml` — auth/login/OAuth endpoints.
-   `backend/services/booking/openapi.yaml` — room search and booking REST API.
-   `backend/services/approval/openapi.yaml` — admin approval HTTP façade (mirrors the gRPC logic).
-   `backend/services/notification/openapi.yaml` — preferences, send, schedule, and history endpoints.
-   Import any file into Swagger UI/ReDoc to explore the contract while coding.
-   Or run the bundled Swagger UI containers (via `docker compose up auth-openapi booking-openapi approval-openapi notification-openapi`) and browse:
    -   Auth: http://localhost:9001
    -   Booking: http://localhost:9002
    -   Approval: http://localhost:9003
    -   Notification: http://localhost:9004

### gRPC via Kong

-   Gateway port `8443` now advertises HTTP/2, so gRPC clients can tunnel through Kong.
-   Approval traffic is exposed only via gRPC; booking and the rest stay on REST/HTTP.
-   Example (using `grpcurl`) to list pending approvals through the gateway:
    ```bash
    grpcurl \
      -insecure \
      -H "Authorization: Bearer <JWT>" \
      localhost:8443 \
      approval.ApprovalService/ListPending
    ```
    Replace `<JWT>` with an admin token issued by the auth service.


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
