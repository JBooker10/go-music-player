CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    "hash" uuid unique not null default gen_random_uuid(),
    first_name VARCHAR (255),
    last_name VARCHAR (255),
    avatar VARCHAR DEFAULT 'https://prd-wret.s3-us-west-2.amazonaws.com/assets/palladium/production/s3fs-public/styles/full_width/public/thumbnails/image/placeholder-profile_0.png',
    dob DATE,
    email TEXT NOT NULL unique,
    password TEXT NOT NULL,
    created_at TIMESTAMP,
    WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)


    DROP TABLE IF EXISTS tracks;
    CREATE TABLE tracks
    (
        id SERIAL PRIMARY KEY,
        "hash" UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
        spotify_id VARCHAR(255) NOT NULL UNIQUE,
        country VARCHAR(255),
        album_id SERIAL,
        album_title VARCHAR(255),
        title VARCHAR(255),
        genre VARCHAR(255),
        duration BIGINT,
        stream TEXT,
        release_date DATE,
        lyrics TEXT,
        collection VARCHAR(255),
        explicit BOOLEAN NOT NULL,
        popularity INT,
        FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE
    );

    CREATE TABLE albums
    (
        id SERIAL PRIMARY KEY,
        "hash" UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
        spotify_id VARCHAR(255) NOT NULL UNIQUE,
        artist_id SERIAL,
        title VARCHAR(255),
        disc_count INT,
        cover VARCHAR(255),
        cover_big VARCHAR(255),
        cover_xl VARCHAR(255),
        release_date DATE,
        track_count INT,
        track_number INT,
        track_list VARCHAR(255),
        record_label VARCHAR(255),
        type VARCHAR(255),
        FOREIGN KEY (artist_id) REFERENCES artists(id)
    );

    CREATE TABLE artists
    (
        id SERIAL PRIMARY KEY,
        "hash" UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
        spotify_id VARCHAR(255) NOT NULL UNIQUE,
        name VARCHAR(255),
        label VARCHAR(255),
        picture VARCHAR(255),
        picture_xl VARCHAR(255)
    );