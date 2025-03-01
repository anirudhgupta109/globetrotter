CREATE DATABASE globetrotter WITH OWNER = postgres;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE table IF NOT EXISTS users (
  	id uuid DEFAULT uuid_generate_v4 () NOT NULL PRIMARY KEY,
  	username VARCHAR(50) NOT NULL UNIQUE,
  	password VARCHAR(100) NOT NULL,
  	auth_token VARCHAR(50)
);