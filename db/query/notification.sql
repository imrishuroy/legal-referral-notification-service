-- name: GetDeviceTokenByUserId :one
SELECT device_token
FROM devices
WHERE user_id = $1;

-- name: CreateNotification :one
INSERT INTO notifications (
    user_id,
    sender_id,
    target_id,
    target_type,
    notification_type,
    message
) VALUES (
     $1, $2, $3, $4, $5, $6
) RETURNING *;
