WITH votes as (
	SELECT
		track_id,
		count(*) as vote_count
		FROM track_voting.votes
		GROUP BY track_id
)
SELECT
	t.id,
	t.url,
	t.filename,
	COALESCE(v.vote_count, 0) as vote_count,
	t.created_by
FROM
	track_voting.tracks t
LEFT JOIN votes v ON t.id = v.track_id
WHERE
	t.has_played = false
order by v.vote_count desc, t.created_by desc
;
