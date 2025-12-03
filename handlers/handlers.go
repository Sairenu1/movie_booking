package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Sairenu1/movie_booking/models"
	"github.com/Sairenu1/movie_booking/store"
)

// Create Booking
func CreateBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var b models.Booking
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// VALIDATION
	if b.ID == "" || b.Movie == "" || b.MovieNumber == "" || b.Seat == "" || b.User == "" {
		http.Error(w, "All fields (id, movie, movie_number, seat, user) are required", http.StatusBadRequest)
		return
	}

	result, err := store.AddBooking(b)
	if err != nil {
		http.Error(w, "Failed to create booking", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

// Get All
func GetBookings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := store.GetAllActive()
	if err != nil {
		http.Error(w, "Failed to fetch bookings", http.StatusInternalServerError)
		return
	}

	if data == nil {
		data = []models.Booking{} // Return empty array instead of null
	}

	json.NewEncoder(w).Encode(data)
}

// Get One
func GetBookingById(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	b, ok := store.GetOneActive(id)
	if !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(b)
}

// Delete
func DeleteBooking(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	ok := store.DeleteBooking(id)
	if !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Update PUT/PATCH
func UpdateBooking(w http.ResponseWriter, r *http.Request) {

	isPatch := r.Method == http.MethodPatch

	id := r.PathValue("id")

	var b models.Booking
	json.NewDecoder(r.Body).Decode(&b)

	result, ok := store.UpdateBooking(id, b, isPatch)
	if !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(result)
}
