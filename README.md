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


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
