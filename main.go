package main

import (
	"log"
	"net/http"

	"github.com/Sairenu1/movie_booking/admin"
	"github.com/Sairenu1/movie_booking/handlers"
	"github.com/Sairenu1/movie_booking/store"
)

// serve login.html as the homepage
func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/login.html")
}

func main() {

	// ---------------- DATABASE ----------------
	if err := store.InitDB(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	router := http.NewServeMux()

	// ---------------- FRONTEND ROUTES ----------------

	// Default route → login page
	router.HandleFunc("/", serveIndex)

	// Serve CSS / JS / Images
	staticFS := http.FileServer(http.Dir("static"))
	router.Handle("/static/", http.StripPrefix("/static/", staticFS))

	// Explicit routes for pages
	router.HandleFunc("/login.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/login.html")
	})

	router.HandleFunc("/dashboard.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/dashboard.html")
	})

	router.HandleFunc("/movies.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/movies.html")
	})

	router.HandleFunc("/bookings.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/bookings.html")
	})

	// ---------------- USER BOOKING APIs ----------------

	router.HandleFunc("GET /bookings", handlers.GetBookings)
	router.HandleFunc("GET /bookings/{id}", handlers.GetBookingById)
	router.HandleFunc("POST /bookings", handlers.CreateBooking)
	router.HandleFunc("PUT /bookings/{id}", handlers.UpdateBooking)
	router.HandleFunc("PATCH /bookings/{id}", handlers.UpdateBooking)
	router.HandleFunc("DELETE /bookings/{id}", handlers.DeleteBooking)

	// ---------------- ADMIN APIs ----------------

	router.HandleFunc("POST /admin/login", admin.AdminLogin)

	// Movie management
	router.HandleFunc("POST /admin/movies", admin.AddMovie)
	router.HandleFunc("GET /admin/movies", admin.GetMovies)
	router.HandleFunc("PUT /admin/movies/{id}", admin.UpdateMovie) // NEW
	router.HandleFunc("PATCH /admin/movies/{id}/deactivate", admin.DeactivateMovie)
	router.HandleFunc("PATCH /admin/movies/{id}/restore", admin.RestoreMovie)

	// Booking management (admin)
	router.HandleFunc("GET /admin/bookings", admin.AdminGetAllBookings)
	router.HandleFunc("PATCH /admin/bookings/{id}/deactivate", admin.DeactivateBooking)
	router.HandleFunc("PATCH /admin/bookings/{id}/restore", admin.RestoreBooking)

	// Dashboard stats
	router.HandleFunc("GET /admin/dashboard/stats", admin.AdminDashboardStats)

	// ---------------- SERVER START ----------------
	log.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", router)
}
