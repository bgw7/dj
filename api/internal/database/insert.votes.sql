INSERT INTO track_voting.votes
(
	track_id,
	user_id
)
VALUES (
	$1,
	$2
	)
ON CONFLICT DO NOTHING
;
