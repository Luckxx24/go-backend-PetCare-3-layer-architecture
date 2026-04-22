-- +goose up

alter table notifications add column is_read bool;


-- +goose down
alter table notifications drop column if EXISTS is_read bool;