
-- +goose up

create type conds as ENUM('makan','sehat','sakit','grooming');

CREATE TABLE pet_status_log (
    ID uuid PRIMARY KEY,
    id_bookings uuid references bookings(ID) ON DELETE CASCADE,
    status conds NOT NULL,
    note TEXT,
    photo_path TEXT,
    created_by uuid references users(ID) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose down
DROP TABLE IF EXISTS pet_status_log;
DROP TYPE IF EXISTS conds;