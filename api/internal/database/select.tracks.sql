WITH votes as (
	SELECT
		url,
		count(*) as vote_count
		FROM track_voting.votes
		GROUP BY url
)
SELECT
	t.id,
	t.url,
	t.filename,
	COALESCE(v.vote_count, 0) as vote_count,
	t.created_by
FROM
	track_voting.tracks t
LEFT JOIN votes v ON t.url = v.url
WHERE
	t.has_played = false
order by v.vote_count desc, t.created_by desc
;
