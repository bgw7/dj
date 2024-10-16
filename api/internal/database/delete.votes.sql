DELETE FROM track_voting.votes
WHERE
	url = $1 AND
	user_id = $2
;
