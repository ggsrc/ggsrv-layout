-- name: AddUser :one
-- -- invalidate: [GetAllUsers]
-- -- timeout: 1s
INSERT INTO public.users (
    name,
    age,
    created_at,
    updated_at
) VALUES (@name, @age, @created_at, @updated_at) RETURNING *;

-- name: GetUserById :one
-- -- cache: 300s
-- -- timeout: 250ms
SELECT * FROM public.users WHERE id = @id;

-- name: GetAllUsers :many
-- -- cache: 30s
-- -- timeout: 500ms
SELECT * FROM public.users ORDER BY id ASC;

-- name: UpdateUserById :one
-- -- invalidate: [GetUserById, GetAllUsers]
-- -- timeout: 1s
UPDATE public.users SET
    name = @name,
    age = @age,
    updated_at = @updated_at
WHERE id = @id RETURNING *;

-- name: DeleteUserById :one
-- -- invalidate: [GetUserById, GetAllUsers]
-- -- timeout: 1s
DELETE FROM public.users WHERE id = @id RETURNING *;

