-- name: CreateMessage :one

insert into message (
     ID,bookings_id,sender_id,receiver_id,message,is_read,created_at
)values (
    $1,$2,$3,$4,$5,$6,$7
)
RETURNING *;

-- name: GetChatInbox :many
SELECT * FROM(
    SELECT DISTINCT ON(bookings_id)
    m.bookings_id,
    m.sender_id,
    m.receiver_id,
    m.message,
    m.is_read,
    m.created_at,
    p.nama as nama_hewan FROM message m INNER JOIN bookings b on bookings_id = b.ID
    INNER JOIN pets p on b.pet_id = p.ID WHERE m.sender_id = $1 OR m.receiver_id =$1 ORDER BY m.bookings_id,m.created_at DESC
)as subquery
ORDER BY created_at DESC
OFFSET $2 LIMIT $3;


-- name: GetHistoryPesan :many
SELECT 
    m.id, m.message, m.is_read, m.created_at, m.sender_id,
    u_sender.nama as pengirim_nama
FROM message m 
INNER JOIN users u_sender on m.sender_id = u_sender.ID
WHERE m.bookings_id = $1
ORDER BY m.created_at ASC
OFFSET $2 LIMIT $3;

-- name: MarkMessageAsRead :exec
update message set is_read = true where bookings_id= $1 and sender_id = $2 and is_read = false;

-- name: DeleteMessage :exec
DELETE from message where ID = $1 and sender_id = $2;

-- name: GetMessagebyIDuser :one

SELECT ID,bookings_id,sender_id,receiver_id,message from message where sender_id = $1; 