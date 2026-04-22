
-- +goose up 
 
 create table notifications(
    ID uuid PRIMARY KEY,
    id_user uuid references users(ID) ON DELETE CASCADE,
    title TEXT,
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose down 

DROP TABLE IF EXISTS notifications;