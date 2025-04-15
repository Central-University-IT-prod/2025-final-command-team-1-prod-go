-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    email varchar(64) primary key,
    username varchar(32),
    "password" text,
    created_at timestamp,
    updated_at timestamp,
    telegram_username text,
    is_admin BOOLEAN,
    UNIQUE(username)
);


CREATE TABLE IF NOT EXISTS places (
    id SERIAL PRIMARY KEY,
    name TEXT,
    description TEXT,
    address TEXT,
    city TEXT
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    images TEXT[],
    user_email TEXT,
    place_id int,
    title TEXT,
    description TEXT,
    genre TEXT,
    author TEXT,
    publication_year INT,
    publisher TEXT,
    condition TEXT,
    -- booked, available, taken
    status TEXT,
    created_at timestamp,
    cover TEXT,
    pages_count int,
    summary TEXT DEFAULT '',

    CONSTRAINT fk_user_email FOREIGN KEY (user_email) REFERENCES users(email)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_place_id FOREIGN KEY (place_id) REFERENCES places(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


CREATE TABLE IF NOT EXISTS favorites (
    id SERIAL PRIMARY KEY,
    user_email varchar(64),
    post_id int,

    CONSTRAINT fk_user_email FOREIGN KEY (user_email) REFERENCES users(email)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_post_id FOREIGN KEY (post_id) REFERENCES posts(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    UNIQUE(user_email, post_id)
);

CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    user_email varchar(64),
    post_id int,
    created_at timestamp,

    CONSTRAINT fk_user_email FOREIGN KEY (user_email) REFERENCES users(email)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_post_id FOREIGN KEY (post_id) REFERENCES posts(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    UNIQUE(user_email, post_id)  
);


CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    target_user_email varchar(64),
    reviewer_user_email varchar(64),
    rating int,
    comment TEXT,
    created_at timestamp,

    CONSTRAINT fk_reviewer_user_email FOREIGN KEY (reviewer_user_email) REFERENCES users(email)
        ON DELETE CASCADE
        ON UPDATE CASCADE,


    CONSTRAINT fk_target_user_email FOREIGN KEY (target_user_email) REFERENCES users(email)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE users;
DROP TABLE posts;
DROP TABLE places;
DROP TABLE reviews;
DROP TABLE bookings;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
