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
EOF