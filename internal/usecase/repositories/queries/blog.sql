-- name: CreateBlog :one
INSERT INTO blog (
    descriptions
) VALUES (
             $1
         )
    RETURNING *;

-- name: GetBlog :one
SELECT * FROM blog
WHERE id = $1 LIMIT 1;

-- name: ListBlog :many
SELECT * FROM blog
ORDER BY created_at
    LIMIT $1
OFFSET $2
;

-- name: DeleteBlog :exec
DELETE FROM blog WHERE id = $1;
