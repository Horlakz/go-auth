# Go Auth API

This is a Go-based authentication API. It provides endpoints for user authentication, authorization, and other related functionalities.

## Getting Started

### Prerequisites

- Go 1.20 or higher
- PostgreSQL v15 and above

### Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/horlakz/go-auth.git
   cd go-auth
   ```

2. Copy the sample environment file and configure it:

   ```sh
   cp .env.sample .env
   ```

3. Install the dependencies:

   ```sh
   go mod download
   ```

### Running the Application

To start the application, run:

```sh
go run main.go
```

To run with binary

```sh
go build -o go-auth
./go-auth
```

## API Endpoints

[postman documentation](https://documenter.getpostman.com/view/26276921/2sAY55ayAf)
