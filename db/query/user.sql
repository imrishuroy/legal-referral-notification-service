-- name: GetUserNameByUserId :one
SELECT first_name, last_name
FROM users
WHERE user_id = $1;
