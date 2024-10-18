INSERT INTO track_voting.tracks
(
	url,
	filename,
	created_by
) VALUES (
	$1,
	$2,
	$3
)
RETURNING id;
