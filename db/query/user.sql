-- name: GetUser :one
select * from users
where id = $1 limit 1;

-- name: CreateUser :one
insert into users(
    id,
    name,
    email,
    password
)values(
    $1,$2,$3,$4
) RETURNING *;

-- name: GetUserByEmail :one
select * from users
where email = $1 limit 1;