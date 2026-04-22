-- +goose up

alter table pet_status_log add column edited_at TIMESTAMPTZ not null;

-- +goose down

alter table pet_status_log DROP column IF EXISTS edited_at TIMESTAMPTZ not null;
