UPDATE track_voting.tracks
SET has_played = $2
WHERE id = $1
;