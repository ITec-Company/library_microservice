CREATE TYPE difficulty AS enum ('junior', 'middle', 'senior');
CREATE TYPE review_source AS enum ('book', 'article', 'video', 'audio');

CREATE TABLE book
(   uuid                    SERIAL PRIMARY KEY,
    title                   CHARACTER VARYING(100) NOT NULL  ,
    direction_uuid          integer NOT NULL ,
    author_uuid             integer NOT NULL ,
    difficulty              difficulty NOT NULL,
    edition_date            DATE NOT NULL ,
    rating                  NUMERIC(9,2) ,
    description             TEXT ,
    url                     TEXT NOT NULL ,
    language                CHARACTER VARYING(15) ,
    tags_uuids              integer [],
    download_count          INTEGER NOT NULL,
    image_url               TEXT DEFAULT '',
    created_at              TIMESTAMP NOT NULL
);

CREATE TABLE article
(   uuid                    SERIAL PRIMARY KEY,
    title                   CHARACTER VARYING(100) NOT NULL  ,
    direction_uuid          integer NOT NULL ,
    author_uuid             integer NOT NULL ,
    difficulty              difficulty NOT NULL,
    edition_date            DATE NOT NULL ,
    rating                  NUMERIC(9,2) ,
    description             TEXT ,
    url                     TEXT NOT NULL ,
    language                CHARACTER VARYING(15) ,
    tags_uuids              integer [],
    download_count          INTEGER NOT NULL,
    image_url               TEXT DEFAULT '',
    created_at              TIMESTAMP NOT NULL
);

CREATE TABLE audio
(   uuid                    SERIAL PRIMARY KEY,
    title                   CHARACTER VARYING(100) NOT NULL  ,
    direction_uuid          integer NOT NULL ,
    difficulty              difficulty NOT NULL,
    rating                  NUMERIC(9,2) ,
    url                     TEXT NOT NULL ,
    language                CHARACTER VARYING(15) ,
    tags_uuids              integer [],
    download_count          INTEGER NOT NULL
);

CREATE TABLE video
(   uuid                    SERIAL PRIMARY KEY,
    title                   CHARACTER VARYING(100) NOT NULL  ,
    direction_uuid          integer NOT NULL ,
    difficulty              difficulty NOT NULL,
    rating                  NUMERIC(9,2) ,
    url                     TEXT NOT NULL ,
    language                CHARACTER VARYING(15) ,
    tags_uuids              integer [],
    download_count          INTEGER NOT NULL
);

CREATE TABLE review
(   uuid                    SERIAL PRIMARY KEY,
    full_name   CHARACTER VARYING(100) NOT NULL  ,
    text                    text NOT NULL  ,
    rating                  NUMERIC(9,2) ,
    source                  review_source NOT NULL ,
    date                    DATE NOT NULL ,
    literature_uuid         integer
);

CREATE TABLE author
(   uuid                    SERIAL PRIMARY KEY,
    full_name               CHARACTER VARYING(100) NOT NULL
);

CREATE TABLE direction
(   uuid                   SERIAL PRIMARY KEY,
    name                   CHARACTER VARYING(100) NOT NULL
);

CREATE TABLE tag
(   uuid                   SERIAL PRIMARY KEY,
    name                   CHARACTER VARYING(100) NOT NULL
);