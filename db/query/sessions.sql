-- name: CreateSession :one
insert into sessions(
    id,
    user_email,
    refresh_token,
    is_blocked,
    expires_at
)values(
    $1,$2,$3,$4,$5
) RETURNING *;

-- name: GetSession :one
select * from sessions where id=$1;