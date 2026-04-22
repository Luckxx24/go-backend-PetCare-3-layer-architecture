

-- name: CreateNewBookings :one

insert into bookings (
    ID,pet_id ,user_id,start_date,end_date,status,created_at
)values(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetBooking :one

SELECT start_date,end_date,status,created_at from bookings where ID = $1;

-- name: GetBookingByStatus :many

SELECT b.start_date,b.end_date,b.status,b.created_at,p.nama as nama_hewan,u.nama as nama_users 
from bookings b INNER JOIN pets p on pet_id = p.ID 
INNER JOIN users u on user_id = u.ID
where b.status = $1
ORDER BY b.start_date ASC
OFFSET $2 LIMIT $3;

-- name: DeleteBooking :exec

DELETE from bookings where ID = $1;

-- name: UpdateBookings :one
update bookings SET start_date = $1,end_date = $2,status = $3 where ID = $4
RETURNING *;

-- name: GetBookingByUserID :one

SELECT ID,pet_id ,user_id,status FROM bookings where ID = $1 and user_id=$2;


