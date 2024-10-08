CREATE SCHEMA IF NOT EXISTS track_voting;
GRANT USAGE ON SCHEMA track_voting TO api_user;
SET search_path TO track_voting;

DROP TABLE IF EXISTS tracks CASCADE;
DROP TABLE IF EXISTS votes CASCADE;

CREATE TABLE IF NOT EXISTS tracks (
  id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  url VARCHAR(500) NOT NULL,
  filename VARCHAR(250),
  has_played BOOL DEFAULT false,
  created_by VARCHAR(100) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS votes (
  track_id integer NOT NULL,
  user_id VARCHAR(100) NOT NULL,
  CONSTRAINT fk_tracks_id
    FOREIGN KEY(track_id)
      REFERENCES tracks (id),
  UNIQUE(track_id, user_id)
);



GRANT ALL ON TABLE track_voting.tracks TO api_user WITH GRANT OPTION;
GRANT ALL ON TABLE track_voting.votes TO api_user WITH GRANT OPTION;
