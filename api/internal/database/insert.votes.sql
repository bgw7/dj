INSERT INTO track_voting.votes
(
	filename,
	url,
	voter_id
)
VALUES (
	$1,
	$2,
	$3
	)
ON CONFLICT DO NOTHING
;
