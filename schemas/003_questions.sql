CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS challenges (
  	id uuid DEFAULT uuid_generate_v4 () NOT NULL PRIMARY KEY,
    inviter VARCHAR(50) REFERENCES users(username) NOT NULL,
    score INTEGER NOT NULL DEFAULT 0,
    correct_answers INTEGER NOT NULL DEFAULT 0,
    incorrect_answers INTEGER NOT NULL DEFAULT 0,
    clues_revealed INTEGER NOT NULL DEFAULT 0,
	is_active BOOLEAN DEFAULT false,
	ended_at TIMESTAMP,
	question_ids UUID[],
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS questions (
  	id uuid DEFAULT uuid_generate_v4 () NOT NULL PRIMARY KEY,
    destination_id UUID REFERENCES destinations(id) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);