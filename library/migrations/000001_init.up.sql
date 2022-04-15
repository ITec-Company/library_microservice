CREATE TYPE diffuculty AS enum ('junior', 'middle', 'senior');

CREATE TABLE book 
(   id                      SERIAL PRIMARY KEY,
    author_id               INTEGER NOT NULL ,
    sub_direction_id        INTEGER NOT NULL ,
    title                   CHARACTER VARYING(100) NOT NULL  ,
    edition_date            DATE NOT NULL ,
    diffuculty diffuculty   NOT NULL ,
    rating                  NUMERIC(9,2) ,
    description             TEXT ,
    language                CHARACTER VARYING(15) ,
    URL                     TEXT NOT NULL ,
    DowloadCount            INTEGER NOT NULL 
);