-- +goose up

alter table pet_status_log alter column created_by set not null;

-- +goose down

alter table pet_status_log alter column created_by drop not null;
