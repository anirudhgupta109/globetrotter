CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS destinations (
  	id uuid DEFAULT uuid_generate_v4 () NOT NULL PRIMARY KEY,
    city VARCHAR(255) NOT NULL UNIQUE,
    country VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS clues (
  	id uuid DEFAULT uuid_generate_v4 () NOT NULL PRIMARY KEY,
    destination_id UUID REFERENCES destinations(id) NOT NULL,
    clue_text TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS fun_facts (
  	id uuid DEFAULT uuid_generate_v4 () NOT NULL PRIMARY KEY,
    destination_id UUID REFERENCES destinations(id) NOT NULL,
    fact_text TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS trivia (
  	id uuid DEFAULT uuid_generate_v4 () NOT NULL PRIMARY KEY,
    destination_id UUID REFERENCES destinations(id) NOT NULL,
    trivia_text TEXT NOT NULL
);

CREATE INDEX ON clues (destination_id);
CREATE INDEX ON destinations (city);