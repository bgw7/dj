INSERT INTO track_voting.tracks
(
	url,
	filename,
	created_by,
	created_with
) VALUES (
	$1,
	$2,
	$3,
	$4
)
RETURNING id, has_played;
