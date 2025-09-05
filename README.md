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


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.