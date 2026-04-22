

-- name: CreateNotifications :one
insert INTO Notifications(
    ID,id_user,title,message,created_at 
)values(
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetNotifications :one
SELECT title,message,created_at from Notifications where ID = $1;

-- name: DeleteNotifications :exec
DELETE from Notifications where ID = $1;

-- name: UpdateNotification :one
Update Notifications set title = $1,message = $2 where ID = $3
RETURNING *;


-- name: GetHistoryNotifications :many
SELECT title,message,created_at from Notifications where id_user = $1
ORDER BY created_at DESC OFFSET $2 LIMIT $3 ;

-- name: MarkNotificationsaAsRead :exec
Update Notifications Set is_read = true where ID = $1 and is_read = false;

-- name: CountUnreadNotifictaions :exec
SELECT COUNT(*) FROM Notifications where is_read = false and id_user =$1;