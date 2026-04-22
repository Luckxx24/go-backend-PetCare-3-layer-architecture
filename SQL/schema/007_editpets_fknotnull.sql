-- +goose up

Alter table pets alter column user_id set not NULL;

-- +goose down 

Alter table pets alter column user_id DROP not NULL;