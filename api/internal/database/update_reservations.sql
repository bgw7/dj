UPDATE la_viajera.reservations
SET 
    start_timestamp = @start_timestamp,
    end_timestamp = @end_timestamp,
    updated_at = current_timestamp,
    updated_by = @updated_by,
    updated_with = current_user
WHERE public_id = @public_id
RETURNING updated_at;