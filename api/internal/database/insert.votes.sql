INSERT INTO track_voting.votes
(
	"url",
	user_id
)
VALUES (
	$1,
	$2
	)
ON CONFLICT DO NOTHING
;
