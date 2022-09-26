package dbrepo

import (
	"context"
	"github.com/revanfall/bookings/internal/models"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `insert into reservations (first_name, last_name, email,phone, 
                          start_date, end_date,room_id, created_at, updated_at) 
						 values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now()).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into the db
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, 
                               created_at, updated_at, restriction_id)
                               values ($1,$2,$3,$4,$5,$6,$7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.CreatedAt,
		r.UpdatedAt,
		r.RestrictionID)
	if err != nil {
		return err
	}
	return nil
}

// SearchAvailabilityByDates returns true if availability exists form roomID, and false if no availability
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select
	count(id) 
	from
		room_restrictions
	where
	    room_id=$1
		$2 < end_date and  $3 >start_date`

	var numRows int
	err := m.DB.QueryRowContext(ctx, query, start, end, roomID).Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `select r.id, r.room_name from rooms r 
				where r.id not in 
				(select room_id from room_restrictions rr 
				where $1 <rr.end_date and $2>rr.start_date)`

	var rooms []models.Room

	rows, err := m.DB.QueryContext(ctx, query,
		start,
		end)

	for rows.Next() {
		var room models.Room
		rows.Scan(&room.ID, &room.RoomName)
		rooms = append(rooms, room)
	}
	if err != nil {
		return rooms, err
	}
	return rooms, nil

}

// GetRoomByID gets a room by id
func (p *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `select id, room_name, created_at, updated_at from rooms where id=$1`

	row := p.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt)
	if err != nil {
		return room, err
	}
	return room, nil
}
