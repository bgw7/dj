CREATE SCHEMA IF NOT EXISTS track_voting;
GRANT USAGE ON SCHEMA track_voting TO api_user;
SET search_path TO track_voting;

DROP TABLE IF EXISTS tracks CASCADE;
DROP TABLE IF EXISTS votes CASCADE;

CREATE TABLE IF NOT EXISTS tracks (
  id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  url VARCHAR(500) NOT NULL UNIQUE,
  filename VARCHAR(250),
  has_played BOOL DEFAULT false,
  created_by VARCHAR(100) NOT NULL,
  created_with VARCHAR(100) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  UNIQUE (url, filename)
);

CREATE TABLE IF NOT EXISTS votes (
  filename VARCHAR(250) NOT NULL,
  url VARCHAR(500) NOT NULL,
  voter_id VARCHAR(100) NOT NULL,
  CONSTRAINT fk_tracks_votes
    FOREIGN KEY(filename, url)
      REFERENCES tracks (filename, url),
  UNIQUE(filename, url, voter_id)
);

GRANT ALL ON TABLE tracks TO api_user WITH GRANT OPTION;
GRANT ALL ON TABLE votes TO api_user WITH GRANT OPTION;
