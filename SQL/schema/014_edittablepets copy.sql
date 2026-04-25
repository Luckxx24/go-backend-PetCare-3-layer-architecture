
-- +goose up

alter table pets add column photo_path text;


-- +goose down 

alter table pets drop column photo_path; 

