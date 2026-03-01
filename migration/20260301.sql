DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE, 
    password VARCHAR(100)
);

INSERT INTO users (username, email, password) VALUES ('testuser', 'testuser@example.com', NULL);

DROP table IF EXISTS media;

CREATE TYPE MediaType AS ENUM ('movie', 'tv');

CREATE TABLE media (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    type MediaType NOT NULL,
    runtime VARCHAR(50) NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    year VARCHAR(4) NOT NULL,
    UNIQUE (title, type, year)
);

INSERT INTO media (title, type, runtime, image_url, year) VALUES
('The Matrix', 'movie', '2h 16m', 'https://image.tmdb.org/t/p/w500/9TGHDvWrqKBzwDxDodHYXEmOE6.jpg', '1999'),
('Fruits Basket (2019)', 'tv', '24m per episode', 'https://image.tmdb.org/t/p/w500/2YbZyWkryXH3sGkMZpS1rFqL9j.jpg', '2019');

DROP table IF EXISTS media_user;

CREATE TYPE Status AS ENUM ('not watched', 'will watch', 'watching', 'watched', 'wont watch');