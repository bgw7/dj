INSERT INTO la_viajera.reservations (
		client_id,
		venue_id,
		start_timestamp,
		end_timestamp,
		created_at,
		created_by,
		created_with
	)
VALUES
	($1, $2, $3, $4, current_timestamp, $5, current_user) RETURNING public_id, created_at;