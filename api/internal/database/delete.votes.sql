DELETE FROM track_voting.votes
WHERE
	filename = $1 AND
	url = $2 AND
	voter_id = $3
;
