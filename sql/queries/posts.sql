-- name: CreatePost :one
INSERT INTO posts (id, title, content, category, tags, createdAt, updatedAt)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    NOW(),
    NOW()
)
RETURNING *;

-- name: GetPosts :many
SELECT * FROM posts;

-- name: GetPost :one
SELECT * FROM posts
WHERE id = $1;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;

-- name: UpdatePost :one
UPDATE posts 
SET updatedAt = NOW(), title = $1, content = $2, category = $3, tags = $4
WHERE id = $5
RETURNING *;

-- name: GetPostsByTerm :many
SELECT * FROM posts
WHERE title ILIKE $1
OR content ILIKE $1
OR category ILIKE $1;