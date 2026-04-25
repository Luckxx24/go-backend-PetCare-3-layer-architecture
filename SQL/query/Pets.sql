

-- name: CreatePets :one

Insert into pets(
    ID,user_id,nama,Jenis,age,created_at,catatan,berat,jenis_kelamin,ras,is_vaxinated,photo_path
) values (
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12
    )
RETURNING *;

-- name: GetPetsListUser :many 
Select nama,jenis,age,ID from pets where user_id = $1 ORDER by created_at DESC OFFSET $2 LIMIT $3;


-- name: GetPetsDetail :one  
Select p.nama,p.jenis,p.age,p.ID,p.catatan,p.berat,p.jenis_kelamin,p.ras,p.is_vaxinated,p.photo_path,u.nama 
from pets p INNER JOIN users U ON U.ID = p.user_id where P.ID = $1;

-- name: GetPetsListSt :many 
SELECT p.nama, p.jenis,u.nama AS owner_name FROM pets p INNER JOIN users u ON p.user_id = u.ID
ORDER BY p.created_at DESC
OFFSET $1 LIMIT $2 ;

-- name: DeletePets :exec

DELETE from pets where ID = $1 and user_id = $2;

-- name: UpdatePets :one
update pets set nama = $1, jenis = $2, age = $3,catatan = $4,berat =$5,jenis_kelamin =$6,ras=$7,is_vaxinated = $8,photo_path = $9 where ID = $10 and user_id = $11
RETURNING *;

-- name: GetPetsByID :one

SELECT id FROM pets WHERE id = $1 AND user_id = $2;
