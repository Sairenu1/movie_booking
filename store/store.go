package store

import (
	"database/sql"
	"log"
	"time"

	"github.com/Sairenu1/movie_booking/models"
)

// Add Booking
func AddBooking(b models.Booking) (models.Booking, error) {

	b.CreatedAt = time.Now()
	b.IsActive = true

	_, err := DB.Exec(`
		INSERT INTO bookings (id, movie, movie_number, seat, user, is_active, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		b.ID, b.Movie, b.MovieNumber, b.Seat, b.User, b.IsActive, b.CreatedAt)

	if err != nil {
		return models.Booking{}, err
	}
	return b, nil
}

// Get All Active Bookings
func GetAllActive() ([]models.Booking, error) {

	rows, err := DB.Query(`
		SELECT id, movie, movie_number, seat, user, is_active, created_at
		FROM bookings
		WHERE is_active = 1`)
	if err != nil {
		return []models.Booking{}, err
	}
	defer rows.Close()

	list := []models.Booking{} // IMPORTANT: initialize an empty slice

	for rows.Next() {
		var b models.Booking
		err := rows.Scan(&b.ID, &b.Movie, &b.MovieNumber, &b.Seat,
			&b.User, &b.IsActive, &b.CreatedAt)
		if err != nil {
			return []models.Booking{}, err
		}
		list = append(list, b)
	}

	return list, nil
}

// Get One Active Booking
func GetOneActive(id string) (models.Booking, bool) {

	var b models.Booking

	err := DB.QueryRow(`
		SELECT id, movie, movie_number, seat, user, is_active, created_at
		FROM bookings
		WHERE id = ? AND is_active = 1`, id).
		Scan(&b.ID, &b.Movie, &b.MovieNumber, &b.Seat,
			&b.User, &b.IsActive, &b.CreatedAt)

	if err == sql.ErrNoRows {
		return models.Booking{}, false
	}
	if err != nil {
		log.Println("GetOne:", err)
		return models.Booking{}, false
	}

	return b, true
}

// Delete Booking (Soft Delete)
func DeleteBooking(id string) bool {

	res, err := DB.Exec(`
		UPDATE bookings SET is_active = 0 WHERE id = ?`, id)

	if err != nil {
		log.Println("Delete:", err)
		return false
	}

	affected, _ := res.RowsAffected()
	return affected > 0
}

// Update Booking
func UpdateBooking(id string, data models.Booking, isPatch bool) (models.Booking, bool) {

	existing, found := GetOneActive(id)
	if !found {
		return models.Booking{}, false
	}

	if isPatch {
		if data.Movie != "" {
			existing.Movie = data.Movie
		}
		if data.MovieNumber != "" {
			existing.MovieNumber = data.MovieNumber
		}
		if data.Seat != "" {
			existing.Seat = data.Seat
		}
		if data.User != "" {
			existing.User = data.User
		}
	} else {
		existing.Movie = data.Movie
		existing.MovieNumber = data.MovieNumber
		existing.Seat = data.Seat
		existing.User = data.User
	}

	_, err := DB.Exec(`
		UPDATE bookings
		SET movie=?, movie_number=?, seat=?, user=?
		WHERE id=? AND is_active=1`,
		existing.Movie, existing.MovieNumber, existing.Seat, existing.User, id)

	if err != nil {
		log.Println("Update:", err)
		return models.Booking{}, false
	}

	return existing, true
}
