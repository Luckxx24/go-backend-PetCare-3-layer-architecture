-- +goose up

create type kelamin as ENUM('Jantan','Betina');

ALter table pets add column Catatan text;
ALter table pets add column Berat Decimal;
ALter table pets add column jenis_kelamin kelamin NOT NULL;
ALter table pets add column ras Varchar;
ALter table pets add column is_vaxinated boolean;

-- +goose down

alter table pets drop if EXISTS Catatan;
alter table pets drop if EXISTS Berat;
alter table pets drop if EXISTS jenis_kelamin;
alter table pets drop type if EXISTS kelamin;
alter table pets drop if EXISTS ras;
alter table pets drop if EXISTS is_vaxinated;



