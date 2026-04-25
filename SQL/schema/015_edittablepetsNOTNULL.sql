
-- +goose up
ALter table pets alter column Catatan SET NOT NULL;
ALter table pets alter column Berat SET NOT NULL;
ALter table pets alter column jenis_kelamin SET NOT NULL;
ALter table pets alter column ras SET NOT NULL;
ALter table pets alter column is_vaxinated SET NOT NULL;
ALter table pets alter column photo_path SET NOT NULL;




-- +goose down 

ALter table pets alter  column Catatan DROP NOT NULL;
ALter table pets alter column Berat DROP  NOT NULL;
ALter table pets alter column jenis_kelamin DROP  NOT NULL;
ALter table pets alter  column ras DROP NOT NULL;
ALter table pets alter column is_vaxinated DROP NOT NULL;
ALter table pets alter column photo_path DROP NOT NULL;

