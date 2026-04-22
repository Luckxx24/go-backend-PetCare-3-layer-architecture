
-- +goose up

create table pets (
    ID uuid PRIMARY KEY,
    user_id uuid references users(ID) ON DELETE CASCADE,
    nama Varchar NOT NULL,
    Jenis Varchar NOT NULL,
    age int ,
    created_at TIMESTAMPTZ  DEFAULT NOW()
);

-- +goose down

DROP TABLE IF EXISTS pets;