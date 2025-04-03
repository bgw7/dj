SELECT
	t.id,
	t.url,
	t.filename
FROM
	track_voting.tracks t
WHERE
	t.id = $1
limit 1
;
