![GGD_Logo](https://github.com/ziadbwdn/borehole_erp/docs/logo/GGD_logo.jpg?raw=true)

# Borehole Data ERP

A robust back-end service built in Go, designed to manage geotechnical and geological data for civil engineering and exploration projects.
This service provides a centralized, secure, and consistent API for handling projects, stations (boreholes), lithology logs, and laboratory test results.


## Core Features

-   **JWT Authentication:** Secure stateless authentication using JSON Web Tokens with an access/refresh token rotation strategy.
-   **Role-Based Access Control (RBAC):** Granular permissions for different user roles (`Engineer`, `Geologist`, `LabTechnician`, `Admin`), ensuring data integrity and security.
-   **Comprehensive Data Management:** Full CRUD (Create, Read, Update, Delete) operations for core geotechnical entities:
    -   Projects
    -   Stations (Boreholes)
    -   Lithology Logs
    -   Laboratory Samples & UCS Test Results
-   **Detailed Audit Trail:** All significant user actions (logins, CRUD operations, data exports) are logged for traceability and security monitoring.
-   **Geospatial Conversion:** Built-in endpoint for converting WGS-84 (Latitude/Longitude) coordinates to UTM for use in GIS and engineering software.
-   **Data Export:** Endpoints to easily export station point lists in both standard geographic (Geo) and UTM formats.
-   **Containerized Environment:** Fully containerized with Docker and Docker Compose for easy setup, development, and deployment.

## Tech Stack

-   **Language:** [Go](https://golang.org/)
-   **Framework:** [Gin Gonic](https://gin-gonic.com/)
-   **ORM:** [GORM](https://gorm.io/)
-   **Database:** [MySQL](https://www.mysql.com/)
-   **Authentication:** [JWT](https://jwt.io/)
-   **Containerization:** [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)

## Getting Started

### Prerequisites

-   [Docker](https://www.docker.com/get-started) and [Docker Compose](https://docs.docker.com/compose/install/)
-   [Go](https://golang.org/dl/) (v1.18 or higher) for local development
-   [Git](https://git-scm.com/)

### Installation & Running

#### Using Docker (Recommended)

This is the simplest way to get the application and its database running.

1.  **Clone the repository:**
    ```bash
    git clone <your-repository-url>
    cd boreholedata-ms
    ```

2.  **Create the configuration file:**
    Copy the example environment file. The default values are configured to work with the `docker-compose.yaml` file.
    ```bash
    cp .env.example .env
    ```

3.  **Run with Docker Compose:**
    This command will build the Go application image, start the MySQL database container, and run the service. The GORM auto-migration will create the necessary tables on startup.
    ```bash
    docker-compose up --build
    ```

The API will be available at `http://localhost:8080`.

#### Running Locally (for Development)

1.  **Clone the repository and navigate into it.**

2.  **Set up the database:**
    Ensure you have a running MySQL instance.

3.  **Create and configure the `.env` file:**
    ```bash
    cp .env.example .env
    ```
    Now, edit the `.env` file and update the `DB_*` variables to match your local MySQL configuration.

4.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

5.  **Run the application:**
    ```bash
    go run cmd/server/main.go
    ```

The API will be available at the port specified in your `.env` file (e.g., `http://localhost:8080`).

## Configuration

The application is configured using environment variables loaded from a `.env` file in the project root.

```ini
# .env.example

# Server Port
PORT=8080

# Database Connection
DB_HOST=db
DB_PORT=3306
DB_USER=user
DB_PASSWORD=password
DB_NAME=boreholedb

# JWT Secret
# IMPORTANT: Use a long, random string for production environments.
JWT_SECRET=your-super-secret-key-that-is-very-long
```

## API Documentation

This project uses Swagger for API documentation. Once the server is running, you can access the interactive Swagger UI at:

**`http://localhost:8080/swagger/index.html`**

## API Endpoints Overview

This is a high-level overview of the available endpoints. For detailed request/response models, please refer to the Swagger documentation.

| Method | Endpoint                             | Description                                  | Required Role(s)                |
| :----- | :----------------------------------- | :------------------------------------------- | :------------------------------ |
| `POST` | `/api/auth/register`                 | Register a new user                          | Public                          |
| `POST` | `/api/auth/login`                    | Authenticate and receive tokens              | Public                          |
| `POST` | `/api/auth/refresh`                  | Refresh token action from login              | Public                          |
| `POST` | `/api/auth/logout`                   | Invalidate the current session               | Public (requires refresh token) |
| `GET`  | `/api/auth/profile`                  | Get the authenticated user's profile         | Any Authenticated User          |
| `PUT`  | `/api/auth/profile`                  | Update the authenticated user's profile      | Any Authenticated User          |
| `POST` | `/api/projects`                      | Create a new project                         | Engineer                        |
| `GET`  | `/api/projects/:id`                  | Get details of a specific project            | Any Authenticated User          |
| `POST` | `/api/stations`                      | Create a new station (borehole)              | Engineer, Geologist             |
| `GET`  | `/api/stations/:id`                  | Get details of a specific station            | Any Authenticated User          |
| `GET`  | `/api/projects/:id/stations/export/geo`  | Export station points as Lat/Lon JSON    | Any Authenticated User          |
| `GET`  | `/api/projects/:id/stations/export/utm`  | Export station points as UTM JSON        | Any Authenticated User          |
| `POST` | `/api/lithology-logs`                | Add a lithology log to a station             | Geologist                       |
| `POST` | `/api/lab-samples`                   | Add a lab sample from a station              | Lab Technician                  |
| `POST` | `/api/ucs-results`                   | Add a UCS test result to a lab sample        | Lab Technician                  |
| `GET`  | `/api/activities`                    | View user activity logs (audit trail)        | Admin                           |


## Project Structure

The project follows a standard Go application layout to enforce separation of concerns.

```
boreholedata-ms/
├── cmd/                # Main application entry points
│   └── server/main.go
├── internal/           # Private application and library code
│   ├── api/            # API layer: handlers, routers, DTOs
│   ├── utils/          # uuid, password, gorm decimal management logic
│   ├── database/       # Database connection and migration
│   ├── interfaces/     # Service and repository contracts (interfaces)
│   ├── models/         # Database models (GORM structs)
│   └── usecases/       # Core business logic: services and repositories
├── pkg/                # Public, reusable libraries (e.g., JWT helpers, converters)
├── Dockerfile          # Docker build file for the Go application
└── docker-compose.yaml # Docker Compose file for development environment
```
## Postman Documentation Link:

Documentation Link: https://documenter.getpostman.com/view/40938916/2sB2xBEVqS
