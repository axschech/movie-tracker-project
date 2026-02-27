DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE, 
    password VARCHAR(100)
);

INSERT INTO users (username, email, password) VALUES ('testuser', 'testuser@example.com', NULL);