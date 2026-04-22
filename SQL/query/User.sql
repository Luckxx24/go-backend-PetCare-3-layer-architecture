-- name: CreateUser :one

Insert into users (
    ID,nama,password,email,role,created_at 
)values(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)

RETURNING *;

-- name: GetUserID :one

Select nama,password,email,role from  users where ID = $1;

-- name: GetUseremail :one

Select ID,nama,password,email,role from  users where email = $1;

-- name: ListsUser :many

Select nama,password,email,role,created_at  from  users order by created_at desc limit $1 offset $2;

-- name: DeleteUser :exec

DELETE from  users where ID = $1;

-- name: UpdateUser :one 

update users set nama = $1, password = $2, email = $3, role = $4 where ID = $5
RETURNING *;

-- name: CountUsers :one
Select COUNT(*) from users;