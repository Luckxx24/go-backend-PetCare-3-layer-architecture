
-- +goose up 


 create type  booking_status as ENUM('PENDING','APPROVED','DONE');


 create TABLE bookings(
    ID uuid PRIMARY KEY,
    pet_id uuid references pets(ID) ON DELETE CASCADE,
    user_id uuid references users(ID) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status booking_status NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose down

DROP TABLE IF EXISTS bookings; 
DROP TYPE IF EXISTS stats;