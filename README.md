# 🚀 Gin GORM CRUD API (User & Todo Management)

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)
[![Gin Framework](https://img.shields.io/badge/Gin-Gonic-00ADD8?logo=gin&logoColor=white)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/GORM-ORM-7FBC39?logo=gopher&logoColor=white)](https://gorm.io/)
[![SQLite](https://img.shields.io/badge/Database-SQLite-003B57?logo=sqlite&logoColor=white)](https://www.sqlite.org/index.html)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## 🌟 Overview

This project is a robust, clean, and well-documented **CRUD (Create, Read, Update, Delete)** RESTful API built with **Go (Golang)**, utilizing the high-performance **Gin** framework and the powerful **GORM** ORM.

It features two primary resources: **Users** and **Todos**, demonstrating a simple **One-to-Many** relationship (one User can have many Todos). The database layer uses a file-based **SQLite** database, making setup incredibly fast.

## 📦 Features

* **RESTful Endpoints:** Full CRUD operations for Users and Todos.
* **Database:** Configured for local **SQLite** for zero-setup development.
* **GORM ORM:** Clean database interactions and auto-migration based on Go structs (Code-First).
* **Swagger Documentation:** Automatically generated OpenAPI 2.0 specification for easy API testing and reference.
* **Structured Handlers:** Logic separated into `handlers` and `models` packages for maintainability.

---

## 🛠️ Installation and Setup

### Prerequisites

You need the following installed on your machine:
* [Go (1.20 or newer)](https://go.dev/dl/)

### Steps

1.  **Clone the Repository:**
    ```bash
    git clone [https://github.com/asirialwis/gin-todo-crud-apis.git](https://github.com/asirialwis/gin-todo-crud-apis.git)
    cd gin-todo-crud-apis
    ```

2.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Install Swagger Generator (if not installed):**
    This tool is required to generate the API documentation.
    ```bash
    go install [github.com/swaggo/swag/cmd/swag@latest](https://github.com/swaggo/swag/cmd/swag@latest)
    ```

4.  **Generate API Documentation:**
    This command reads the annotations in your handler functions and generates the `docs/swagger.json` file.
    ```bash
    swag init -g main.go
    ```

5.  **Run the Application:**
    The application will automatically connect to SQLite and run GORM migrations to create the `users` and `todos` tables if they don't exist.
    ```bash
    go run main.go
    ```
    The server will start at `http://localhost:8080`.

---

## 📝 API Documentation (Swagger UI)

Once the server is running, the full API documentation is available through the Swagger UI:

🔗 **Swagger UI Link:** `http://localhost:8080/swagger/index.html`

This interface allows you to view, test, and interact with all endpoints, grouped neatly into **Users** and **Todos** categories.

---

## 🎯 API Endpoints

The base URL for the API is `http://localhost:8080/`.

### User Endpoints (`/users`)

| Method | Path | Description |
| :--- | :--- | :--- |
| `POST` | `/users` | Create a new user. |
| `GET` | `/users` | Retrieve all users (with associated todos preloaded). |
| `GET` | `/users/:id` | Retrieve a single user by ID. |
| `PATCH` | `/users/:id` | Update a user's details. |
| `DELETE`| `/users/:id` | Soft-delete a user (keeps record, sets `DeletedAt`). |

### Todo Endpoints (`/todos`)

| Method | Path | Description |
| :--- | :--- | :--- |
| `POST` | `/todos` | Create a new todo item (requires existing `user_id`). |
| `GET` | `/todos` | Retrieve all todo items. |
| `GET` | `/todos/:id` | Retrieve a single todo by ID. |
| `PATCH` | `/todos/:id` | Update a todo item (e.g., mark as completed). |
| `DELETE`| `/todos/:id` | Soft-delete a todo item. |

---

## 📂 Project Structure

A clean project structure for maintainability:
