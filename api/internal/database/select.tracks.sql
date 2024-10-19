WITH votes as (
	SELECT
		filename,
		url,
		count(*) as vote_count
		FROM track_voting.votes
		GROUP BY filename, url
)
SELECT
	t.id,
	t.url,
	t.filename,
	COALESCE(v.vote_count, 0) as vote_count,
	t.created_by
FROM
	track_voting.tracks t
LEFT JOIN votes v 
	ON t.filename = v.filename AND
	   t.url = v.url
WHERE
	t.has_played = false
order by v.vote_count desc, t.created_by desc
;
