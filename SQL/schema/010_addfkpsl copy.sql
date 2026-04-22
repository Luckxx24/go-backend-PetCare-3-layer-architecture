-- +goose up

alter table pet_status_log alter column created_at set not null;

-- +goose down

alter table pet_status_log alter column created_at drop not null;
