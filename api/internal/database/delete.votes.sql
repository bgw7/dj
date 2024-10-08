DELETE FROM track_voting.votes
WHERE
	track_id = $1 AND
	user_id = $2
;
