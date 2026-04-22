
-- +goose up

create table message (
    ID uuid PRIMARY KEY,
    bookings_id uuid references bookings(ID) ON DELETE CASCADE,
    sender_id uuid references users(ID) ON DELETE CASCADE,
    receiver_id uuid references users(ID) ON DELETE CASCADE,
    message TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose down

DROP TABLE IF EXISTS message;