package admin

import (
	"encoding/json"
	"net/http"

	"github.com/Sairenu1/movie_booking/models"
	"github.com/Sairenu1/movie_booking/store"
)

// ---------------- ADMIN LOGIN ------------------

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if ValidateAdmin(req.Username, req.Password) {
		w.Write([]byte("Login successful"))
		return
	}

	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

// ---------------- MOVIE CRUD (SOFT DELETE + UPDATE) --------------------

type Movie struct {
	Title       string `json:"title"`
	MovieNumber string `json:"movie_number"`
	Genre       string `json:"genre"`
	Duration    int    `json:"duration"`
}

// ---------------- ADD MOVIE ----------------

func AddMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var m Movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if m.Title == "" || m.MovieNumber == "" {
		http.Error(w, "title and movie_number are required", http.StatusBadRequest)
		return
	}

	_, err := store.DB.Exec(`
		INSERT INTO movies (title, movie_number, genre, duration, is_active)
		VALUES (?, ?, ?, ?, 1)
	`, m.Title, m.MovieNumber, m.Genre, m.Duration)

	if err != nil {
		http.Error(w, "Failed to add movie", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Movie added successfully",
	})
}

// ---------------- GET ACTIVE MOVIES ----------------

func GetMovies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := store.DB.Query(`
		SELECT id, title, movie_number, genre, duration
		FROM movies
		WHERE is_active = 1
	`)
	if err != nil {
		http.Error(w, "Failed to fetch movies", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []map[string]interface{}

	for rows.Next() {
		var id, duration int
		var title, movieNumber, genre string

		if err := rows.Scan(&id, &title, &movieNumber, &genre, &duration); err != nil {
			http.Error(w, "Failed to read movie", http.StatusInternalServerError)
			return
		}

		movies = append(movies, map[string]interface{}{
			"id":           id,
			"title":        title,
			"movie_number": movieNumber,
			"genre":        genre,
			"duration":     duration,
		})
	}

	json.NewEncoder(w).Encode(movies)
}

// ---------------- UPDATE MOVIE (PUT) ----------------

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")

	var m Movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if m.Title == "" || m.MovieNumber == "" {
		http.Error(w, "title and movie_number are required", http.StatusBadRequest)
		return
	}

	res, err := store.DB.Exec(`
		UPDATE movies
		SET title = ?, movie_number = ?, genre = ?, duration = ?
		WHERE id = ? AND is_active = 1
	`, m.Title, m.MovieNumber, m.Genre, m.Duration, id)

	if err != nil {
		http.Error(w, "Failed to update movie", http.StatusInternalServerError)
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		http.Error(w, "Movie not found or inactive", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Movie updated successfully",
	})
}

// ---------------- SOFT DELETE MOVIE ----------------

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")

	res, err := store.DB.Exec(`
		UPDATE movies
		SET is_active = 0
		WHERE id = ?
	`, id)

	if err != nil {
		http.Error(w, "Failed to delete movie", http.StatusInternalServerError)
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ================= MOVIE SOFT DELETE =================

// Soft delete movie
func DeactivateMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	_, err := store.DB.Exec(`UPDATE movies SET is_active = 0 WHERE id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to deactivate movie", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Movie deactivated"))
}

// Restore movie
func RestoreMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	_, err := store.DB.Exec(`UPDATE movies SET is_active = 1 WHERE id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to restore movie", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Movie restored"))
}

// --------------- BOOKING MGMT -------------------

func AdminGetAllBookings(w http.ResponseWriter, r *http.Request) {
	rows, err := store.DB.Query(`SELECT * FROM bookings`)
	if err != nil {
		http.Error(w, "Failed to fetch bookings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bookings []models.Booking

	for rows.Next() {
		var b models.Booking
		rows.Scan(&b.ID, &b.Movie, &b.MovieNumber, &b.Seat, &b.User, &b.IsActive, &b.CreatedAt)
		bookings = append(bookings, b)
	}

	json.NewEncoder(w).Encode(bookings)
}

// --------------- DEACTIVATE BOOKINGS -------------------

func DeactivateBooking(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	store.DB.Exec(`UPDATE bookings SET is_active=0 WHERE id=?`, id)

	w.Write([]byte("Booking deactivated"))
}

// --------------- RESTORE BOOKINGS -------------------

func RestoreBooking(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	store.DB.Exec(`UPDATE bookings SET is_active=1 WHERE id=?`, id)

	w.Write([]byte("Booking restored"))
}

// --------------- DASHBOARD STATS ----------------

func AdminDashboardStats(w http.ResponseWriter, r *http.Request) {

	var total, active, inactive int

	store.DB.QueryRow(`SELECT COUNT(*) FROM bookings`).Scan(&total)
	store.DB.QueryRow(`SELECT COUNT(*) FROM bookings WHERE is_active=1`).Scan(&active)
	store.DB.QueryRow(`SELECT COUNT(*) FROM bookings WHERE is_active=0`).Scan(&inactive)

	result := map[string]int{
		"total_bookings":    total,
		"active_bookings":   active,
		"inactive_bookings": inactive,
	}

	json.NewEncoder(w).Encode(result)
}
