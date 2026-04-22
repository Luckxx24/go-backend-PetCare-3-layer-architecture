
-- +goose up

create type  role as ENUM('User','Staff','Admin');

 create table users(
    ID uuid PRIMARY KEY,
    nama Varchar NOT NULL,
    password TEXT NOT NULL,
    email TEXT NOT NULL,
    role role NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
 );

-- +goose down

 DROP table IF EXISTS users;
 DROP TYPE IF EXISTS role;