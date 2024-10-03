-- name: GetDeviceTokenByUserId :one
SELECT device_token
FROM devices
WHERE user_id = $1;
