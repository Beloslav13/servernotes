#!/bin/bash
whoami
psql << EOF
ALTER USER postgres WITH PASSWORD '123';
CREATE USER admin WITH PASSWORD 'devpass';
CREATE DATABASE servernotes_db;
GRANT ALL PRIVILEGES ON DATABASE servernotes_db TO admin;
ALTER ROLE admin SUPERUSER CREATEDB;
\c servernotes_db;

CREATE TABLE IF NOT EXISTS persons(
    id serial primary key,
    tg_user_id bigint not null UNIQUE,
    username varchar(255) not null UNIQUE,
    created timestamp default current_timestamp
);
INSERT INTO persons(tg_user_id, username) VALUES (128, 'Ivan');


CREATE TABLE categories(
    id serial primary key,
    person_id int,
    FOREIGN KEY (person_id) REFERENCES persons(id) ON DELETE CASCADE,
    name varchar(255) not null UNIQUE,
    created timestamp default current_timestamp
);
INSERT INTO categories (person_id, name) VALUES (1, 'Другое');

CREATE TABLE notes(
    id serial primary key,
    person_id int,
    FOREIGN KEY (person_id) REFERENCES persons(id) ON DELETE CASCADE,
    category_id int,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    name text not null,
    created timestamp default current_timestamp
);
INSERT INTO notes (person_id, category_id, name) VALUES (1, 1, 'https://google.com');

EOF