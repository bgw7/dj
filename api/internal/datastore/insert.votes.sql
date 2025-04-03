INSERT INTO track_voting.votes
(
	track_id,
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
