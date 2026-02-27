DROP table IF EXISTS media;

CREATE TYPE MediaType AS ENUM ('movie', 'tv');

CREATE TABLE media (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    type MediaType NOT NULL,
    runtime VARCHAR(50) NOT NULL,
    image_url VARCHAR(255) NOT NULL
);

DROP table IF EXISTS media_user;

CREATE TYPE Status AS ENUM ('not watched', 'will watch', 'watching', 'watched', 'wont watch');

CREATE TABLE media_user (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    media_id INT NOT NULL,
    status Status NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (media_id) REFERENCES media(id)
);

