# Movie Booking System (Go + MySQL + HTML/CSS/JS)

This project is a simple full-stack Movie Booking System built using Golang for the backend, MySQL as the database, and HTML/CSS/JavaScript for the admin dashboard interface.

The system allows users to create and manage movie bookings, while the admin dashboard provides complete control over movies, bookings, and statistics.

## Overview

- Backend built using **Golang (net/http)** with modular handlers, models, and store layers.
- **MySQL** is used to store bookings, movies, and admin users.
- The admin dashboard is built using **pure HTML, CSS, and JavaScript**, and is served directly from the Go backend.
- The system supports:
  - Creating bookings
  - Updating bookings (PUT/PATCH)
  - Soft deleting bookings
  - Admin login
  - Managing movies (add, update, delete, restore)
  - Viewing all bookings
  - Dashboard statistics showing total, active, and inactive bookings

## Features

- User Booking APIs: create, update, list, get by ID, and soft delete.
- Admin APIs: login, movie management, booking management, and dashboard stats.
- Lightweight frontend with HTML/CSS/JS.
- Static frontend files served directly from the Go server.
- Clean folder structure that separates backend logic, database layer, and frontend UI.

## How to Run

1. Create the MySQL database and required tables.
2. Update MySQL credentials in `store/db.go`.
3. Install dependencies using: go mod tidy
4. Start the server: go run .
5. Open the admin dashboard: http://localhost:8080/login.html

This project is designed to be simple, easy to understand, and ideal for learning.
