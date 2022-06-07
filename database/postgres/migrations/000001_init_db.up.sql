CREATE TYPE difficulty AS enum ('junior', 'middle', 'senior', '');

CREATE TYPE review_source AS enum ('book', 'article', 'video', 'audio', '');

CREATE TABLE article
(   uuid                    SERIAL PRIMARY KEY,
    title                   CHARACTER VARYING(100) NOT NULL ,
    direction_uuid          integer NOT NULL,
    author_uuid             integer NOT NULL,
    difficulty              difficulty NOT NULL,
    edition_date            DATE NOT NULL,
    rating                  NUMERIC(9,2) NOT NULL DEFAULT 0.0,
    all_grades              NUMERIC(9,2)[],
    description             TEXT,
    local_url               TEXT DEFAULT '',
    image_url               TEXT DEFAULT '',
    web_url                 TEXT DEFAULT '',
    language                CHARACTER VARYING(15),
    tags_uuids              integer [],
    download_count          INTEGER NOT NULL DEFAULT 0,
    created_at              TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE TABLE audio
(   uuid                    SERIAL PRIMARY KEY,
    title                   CHARACTER VARYING(100) NOT NULL ,
    direction_uuid          integer NOT NULL,
    difficulty              difficulty NOT NULL,
    rating                  NUMERIC(9,2) NOT NULL DEFAULT 0.0,
    all_grades              NUMERIC(9,2)[],
    local_url               TEXT NOT NULL,
    language                CHARACTER VARYING(15),
    tags_uuids              integer [],
    download_count          INTEGER NOT NULL DEFAULT 0,
    created_at              TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE TABLE book
(   uuid                    SERIAL PRIMARY KEY,
    title                   CHARACTER VARYING(100) NOT NULL,
    direction_uuid          integer NOT NULL,
    author_uuid             integer NOT NULL,
    difficulty              difficulty NOT NULL,
    edition_date            DATE NOT NULL,
    rating                  NUMERIC(9,2) NOT NULL DEFAULT 0.0,
    all_grades              NUMERIC(9,2)[],
    description             TEXT,
    local_url               TEXT DEFAULT '',
    image_url               TEXT DEFAULT '',
    language                CHARACTER VARYING(15),
    tags_uuids              integer [],
    download_count          INTEGER NOT NULL DEFAULT 0,
    created_at              TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE TABLE video
(   uuid                    SERIAL PRIMARY KEY,
    title                   CHARACTER VARYING(100) NOT NULL ,
    direction_uuid          integer NOT NULL,
    difficulty              difficulty NOT NULL,
    rating                  NUMERIC(9,2) NOT NULL DEFAULT 0.0,
    all_grades              NUMERIC(9,2)[],
    local_url               TEXT DEFAULT '',
    web_url                 TEXT DEFAULT '',
    language                CHARACTER VARYING(15),
    tags_uuids              integer [],
    download_count          INTEGER NOT NULL DEFAULT 0,
    created_at              TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE TABLE review
(   uuid                    SERIAL PRIMARY KEY,
    full_name               CHARACTER VARYING(100) NOT NULL ,
    text                    text NOT NULL ,
    rating                  NUMERIC(9,2) NOT NULL DEFAULT 0.0,
    all_grades              NUMERIC(9,2)[],
    source                  review_source NOT NULL,
    date                    DATE NOT NULL,
    literature_uuid         integer,
    created_at              TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE TABLE author
(   uuid                    SERIAL PRIMARY KEY,
    full_name               CHARACTER VARYING(100) NOT NULL,
    created_at              TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE TABLE direction
(   uuid                   SERIAL PRIMARY KEY,
    name                   CHARACTER VARYING(100) NOT NULL,
    created_at             TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

CREATE TABLE tag
(   uuid                   SERIAL PRIMARY KEY,
    name                   CHARACTER VARYING(100) NOT NULL,
    created_at             TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);