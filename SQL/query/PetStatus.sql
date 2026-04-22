

-- name: CreateNewLog :one

Insert into pet_status_log(
    ID,id_bookings,status,note,photo_path,created_by,created_at,edited_at
) values (
    $1,$2,$3,$4,$5,$6,$7,$8
)
RETURNING *;

-- name: GetLog :one

SELECT psl.status,psl.note,psl.photo_path,psl.created_by,psl.created_at,psl.edited_at,p.nama as nama_hewan, p.jenis as jenis_hewan from pet_status_log psl 
INNER JOIN bookings b on psl.id_bookings = b.ID 
INNER JOIN pets p on b.pet_id = p.ID
where psl.id_bookings = $1;

-- name: GetAllLog :many
SELECT ps.status,ps.note,ps.photo_path,ps.created_by,ps.created_at,ps.edited_at,p.nama as nama_hewan, u.nama as nama_users from pet_status_log ps 
INNER JOIN bookings b on ps.id_bookings = b.ID 
INNER JOIN pets p on b.pet_id = p.ID
INNER JOIN users u on b.user_id = u.ID
where u.ID = $1
ORDER BY ps.created_at DESC
OFFSET $2 LIMIT $3;

-- name: DeleteLog :exec
DELETE from pet_status_log where ID = $1;

-- name: UpdateLog :one 
update pet_status_log set status = $1, note = $2, photo_path = $3,edited_at = $4 where ID = $5
RETURNING *; 

-- name: GetLOgbyIDbooking :one
SELECT id_bookings,status FROM pet_status_log where id_bookings = $1;