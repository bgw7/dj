INSERT (
	track_id,
	user_id
)
INTO track_voting.votes 
VALUES (
	$1,
	$2
	)
ON CONFLICT DO NOTHING
;
