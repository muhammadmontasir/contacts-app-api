# Contact Application API

Welcome to the **Contact Application API**! This API allows users to register, activate their accounts, authenticate, and perform CRUD (Create, Read, Update, Delete) operations on their contacts. Authentication is handled using JSON Web Tokens (JWT).

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [Project Initialization](#project-initialization)
  - [Prerequisites](#prerequisites)
  - [Clone the Repository](#clone-the-repository)
  - [Configure Environment Variables](#configure-environment-variables)
  - [Database Setup with Docker](#database-setup-with-docker)
  - [Run Database Migrations](#run-database-migrations)
  - [Start the Server](#start-the-server)
- [API Documentation](#api-documentation)
  - [1. User Registration](#1-user-registration)
  - [2. User Activation](#2-user-activation)
  - [3. User Authentication](#3-user-authentication)
  - [4. Create Contact](#4-create-contact)
  - [5. Get All Contacts](#5-get-all-contacts)
  - [6. Get Contact Details](#6-get-contact-details)
  - [7. Update Contact](#7-update-contact)
  - [8. Delete Contact](#8-delete-contact)
- [Error Handling](#error-handling)

---

## Features

- **User Registration:** Create a new user account.
- **User Activation:** Activate user accounts via token.
- **User Authentication:** Authenticate users and provide JWT tokens.
- **Contact Management:** Create, read, update, and delete contacts associated with authenticated users.
- **Pagination:** Retrieve contacts with pagination support.
- **JWT-Based Authentication:** Secure endpoints using JWT tokens.

## Technologies Used

- **Go:** Programming language.
- **Gorilla Mux:** HTTP router and URL matcher.
- **GORM:** ORM library for database interactions.
- **PostgreSQL:** Database management system.
- **Docker & Docker Compose:** Containerization and orchestration for the database.
- **JWT:** JSON Web Tokens for authentication.
- **Go Modules:** Dependency management.

---

## Project Initialization

Follow the steps below to set up and run the Contact Application API on your local machine.

### Prerequisites

Ensure you have the following installed on your system:

- **Go:** [Download and Install Go](https://golang.org/dl/)
- **Git:** [Download and Install Git](https://git-scm.com/downloads)
- **Docker & Docker Compose:** [Download and Install Docker](https://www.docker.com/get-started)
- **Environment Manager (Optional):** Tools like `direnv` or `.env` files for managing environment variables.

### Clone the Repository

```bash
git clone https://github.com/muhammadmontasir/contacts-app-api
cd contacts-app-api
```

### Configure Environment Variables

Create a `.env` file in the root directory of the project to store your environment variables. Below is a sample configuration:

```env
# Server Configuration
PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=54322
DB_USER=postgres
DB_PASSWORD=
DB_NAME=postgres

# JWT Configuration
JWT_SECRET=
```

**Note:** Replace the placeholder values (`your_jwt_secret_key`) with your actual secret keys. Ensure that the `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, and `DB_NAME` match the settings in the provided `docker-compose.yml`.

### Database Setup with Docker

The project uses PostgreSQL as its database, which is containerized using Docker Compose. Follow these steps to set up the database:

1. **Ensure Docker is Running:**

   Make sure Docker is installed and the Docker daemon is running on your machine.

   ```bash
   docker --version
   docker compose --version
   ```

2. **Start the PostgreSQL Container:**

   Run the following command to start the PostgreSQL service:

   ```bash
   docker compose up -d
   ```

   - The `-d` flag runs the container in detached mode.
   - Verify that the container is running:

     ```bash
     docker ps
     ```

     You should see `contacts-app-api-postgres` listed among the running containers.

3. **Verify Database Connection:**

   Ensure that the application can connect to the PostgreSQL database using the provided credentials.

   - **Host:** `localhost`
   - **Port:** `54322`
   - **User:** `your_db_user`
   - **Password:** `your_password`
   - **Database Name:** `postgres`

   > **Note:** If you need to change these settings, update both the `docker-compose.yml` and the `.env` file accordingly.

### Run Database Migrations

The project uses GORM for ORM. Upon running the application, GORM will automatically migrate the schema based on the models defined in the code.

**Alternatively**, if you have migration scripts or want to manage migrations manually, follow your migration strategy here.

### Start the Server

Build and run the application using the following commands:

```bash
go build -o contacts-app-api
./contacts-app-api
```

**Or**, run directly using `go run` (useful during development):

```bash
go run cmd/server/main.go
```

The server should start and listen on the port specified in your `.env` file (default is `8080`).

**Access the API:** [http://localhost:8080](http://localhost:8080)

---

## API Documentation

Below are the available API endpoints along with their usage details.

### 1. User Registration

Register a new user account.

- **Endpoint:** `POST /api/v1/users`
- **Description:** Creates a new user with the provided username, email, and password.
- **Request:**

  ```http
  POST http://localhost:8080/api/v1/users
  Content-Type: application/json

  {
    "username": "testuser",
    "email": "testuser@example.com",
    "password": "securepassword123"
  }
  ```

- **Response:**

  - **Status:** `201 Created`
  - **Body:**

    ```json
    {
      "activation_token": "token_received_in_email"
    }
    ```

- **Notes:**
  - After registration, an activation token is sent to the user's email for account activation.

### 2. User Activation

Activate a registered user account using the activation token.

- **Endpoint:** `POST /api/v1/users/activate`
- **Description:** Activates the user account associated with the provided email and activation token.
- **Request:**

  ```http
  POST http://localhost:8080/api/v1/users/activate
  Content-Type: application/json

  {
    "email": "testuser@example.com",
    "activation_token": "token_received_in_email"
  }
  ```

- **Response:**

  - **Status:** `200 OK`
  - **Body:**

    ```json
    {
      "message": "User activated successfully"
    }
    ```

### 3. User Authentication

Authenticate a user and receive a JWT token for authorized requests.

- **Endpoint:** `POST /api/v1/token/auth`
- **Description:** Authenticates the user using email and password, returning a JWT token upon successful authentication.
- **Request:**

  ```http
  POST http://localhost:8080/api/v1/token/auth
  Content-Type: application/json

  {
    "email": "testuser@example.com",
    "password": "securepassword123"
  }
  ```

- **Response:**

  - **Status:** `200 OK`
  - **Body:**

    ```json
    {
      "token": "your_generated_jwt_token_here"
    }
    ```

- **Notes:**
  - Save the returned JWT token for authenticating subsequent requests.

### 4. Create Contact

Create a new contact associated with the authenticated user.

- **Endpoint:** `POST /api/v1/contacts`
- **Description:** Creates a new contact with the provided name, email, and phone number.
- **Headers:**

  ```
  Authorization: Bearer <your_jwt_token>
  ```

- **Request:**

  ```http
  POST http://localhost:8080/api/v1/contacts
  Content-Type: application/json
  Authorization: Bearer <your_jwt_token>

  {
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "1234567890"
  }
  ```

- **Response:**

  - **Status:** `201 Created`
  - **Body:**

    ```json
    {
      "id": 1,
      "user_id": 123, // Example user ID
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "1234567890",
      "created_at": "2023-10-10T12:34:56Z",
      "updated_at": "2023-10-10T12:34:56Z"
    }
    ```

### 5. Get All Contacts

Retrieve all contacts associated with the authenticated user, with pagination support.

- **Endpoint:** `GET /api/v1/contacts`
- **Description:** Retrieves a paginated list of contacts for the authenticated user.
- **Headers:**

  ```
  Authorization: Bearer <your_jwt_token>
  ```

- **Query Parameters:**

  - `page` (optional): Page number (default: 1)
  - `page_size` (optional): Number of contacts per page (default: 10)

- **Request:**

  ```http
  GET http://localhost:8080/api/v1/contacts?page=1&page_size=10
  Authorization: Bearer <your_jwt_token>
  ```

- **Response:**

  - **Status:** `200 OK`
  - **Body:**

    ```json
    {
      "contacts": [
        {
          "id": 1,
          "user_id": 123,
          "name": "John Doe",
          "email": "john@example.com",
          "phone": "1234567890",
          "created_at": "2023-10-10T12:34:56Z",
          "updated_at": "2023-10-10T12:34:56Z"
        }
        // More contacts...
      ],
      "total": 25,
      "page": 1,
      "pageSize": 10
    }
    ```

### 6. Get Contact Details

Retrieve details of a specific contact by its ID.

- **Endpoint:** `GET /api/v1/contacts/{contact_id}`
- **Description:** Retrieves detailed information about a specific contact belonging to the authenticated user.
- **Headers:**

  ```
  Authorization: Bearer <your_jwt_token>
  ```

- **Path Parameters:**

  - `contact_id`: The ID of the contact to retrieve.

- **Request:**

  ```http
  GET http://localhost:8080/api/v1/contacts/1
  Authorization: Bearer <your_jwt_token>
  ```

- **Response:**

  - **Status:** `200 OK`
  - **Body:**

    ```json
    {
      "id": 1,
      "user_id": 123,
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "1234567890",
      "created_at": "2023-10-10T12:34:56Z",
      "updated_at": "2023-10-10T12:34:56Z"
    }
    ```

### 7. Update Contact

Update details of an existing contact.

- **Endpoint:** `PATCH /api/v1/contacts/{contact_id}`
- **Description:** Updates the specified contact's information.
- **Headers:**

  ```
  Authorization: Bearer <your_jwt_token>
  ```

- **Path Parameters:**

  - `contact_id`: The ID of the contact to update.

- **Request:**

  ```http
  PATCH http://localhost:8080/api/v1/contacts/1
  Content-Type: application/json
  Authorization: Bearer <your_jwt_token>

  {
    "name": "John Updated",
    "email": "johnupdated@example.com"
  }
  ```

- **Response:**

  - **Status:** `200 OK`
  - **Body:**

    ```json
    {
      "id": 1,
      "user_id": 123,
      "name": "John Updated",
      "email": "johnupdated@example.com",
      "phone": "1234567890",
      "created_at": "2023-10-10T12:34:56Z",
      "updated_at": "2023-10-10T13:00:00Z"
    }
    ```

### 8. Delete Contact

Delete a specific contact by its ID.

- **Endpoint:** `DELETE /api/v1/contacts/{contact_id}`
- **Description:** Deletes the specified contact belonging to the authenticated user.
- **Headers:**

  ```
  Authorization: Bearer <your_jwt_token>
  ```

- **Path Parameters:**

  - `contact_id`: The ID of the contact to delete.

- **Request:**

  ```http
  DELETE http://localhost:8080/api/v1/contacts/1
  Authorization: Bearer <your_jwt_token>
  ```

- **Response:**

  - **Status:** `204 No Content`
  - **Body:** *(Empty)*

---

## Error Handling

The API uses standard HTTP status codes to indicate success or failure of API requests. Below are common status codes you may encounter:

- **200 OK:** The request was successful.
- **201 Created:** The resource was successfully created.
- **204 No Content:** The resource was successfully deleted.
- **400 Bad Request:** The request was invalid or cannot be served.
- **401 Unauthorized:** Authentication failed or user does not have permissions.
- **404 Not Found:** The requested resource could not be found.
- **500 Internal Server Error:** An error occurred on the server.

**Example Error Response:**

```json
{
  "message": "Invalid email or password"
}
```

---

## Contact

For any inquiries or support, please contact [montasircste@gmail.com](mailto:montasircste@gmail.com).
