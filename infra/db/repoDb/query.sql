----------------------------------------------------------------------------------------------------------------------------

-- name: SearchByTitleOrDescription :many
SELECT url, title
FROM pages
WHERE title ILIKE '%' || $1 || '%' -- Case-insensitive search in title
   OR description ILIKE '%' || $1 || '%';
-- Case-insensitive search in description

-- name: SearchByContent :many
SELECT url, title, (words ->> $1)::int AS frequency -- Use ->> for text extraction
FROM pages
WHERE words ? $1 -- Check if the word exists in the words JSONB column
ORDER BY frequency DESC NULLS LAST;

-- name: SearchBy :many
SELECT DISTINCT url, title
FROM pages
WHERE (words ->> $1)::int IS NOT NULL
   OR title ILIKE '%' || $1 || '%'
   OR description ILIKE '%' || $1 || '%'
ORDER BY url;

-- name: GetPage :one
SELECT *
FROM pages
WHERE id = $1 LIMIT 1;

-- name: GetPageByUrl :one
SELECT *
FROM pages
WHERE url = $1 LIMIT 1;

-- name: InsertPage :one
INSERT INTO pages (url, title, description, meta, visited, words, links)
VALUES (@url::text, @title::text, @description::text, @meta::jsonb, @visited::bool, @words::jsonb, @links::text[]) RETURNING *;

-- name: UpdatePage :one
UPDATE pages
SET title       = $2,
    description = $3,
    meta        = $4,
    visited     = $5,
    timestamp   = $6,
    words       = $7
WHERE id = $1 RETURNING *;

-- name: DeletePage :one
DELETE
FROM pages
WHERE id = $1 RETURNING *;

-- name: DeletePageByUrl :one
DELETE
FROM pages
WHERE url = $1 RETURNING *;

-- name: GetPages :many
SELECT *
FROM pages limit $1
offset $2;

-- name: GetPageByDate :many
SELECT *
FROM pages
WHERE timestamp >= $1 AND timestamp <= $2;


----------------------------------------------------------------------------------------------------------------------------


-- name: QueueInsert :one
INSERT INTO queue (url, depth, status)
VALUES (@url::text, @depth::int, @status::int) RETURNING *;

-- name: QueueMarkAs :exec
UPDATE queue
SET status = @status::int
WHERE id = @id::uuid;

-- name: QueueMarkAsByUrl :exec
UPDATE queue
SET status = @status::int
WHERE url = @url::text;

-- name: QueueGetNext :one
SELECT *
FROM queue
WHERE status = 1
ORDER BY timestamp ASC LIMIT 1;

-- name: QueueGetAll :many
SELECT *
FROM queue
ORDER BY timestamp ASC;


-- name: QueueGetAllByStatus :many
SELECT *
FROM queue
WHERE status = @status::int
ORDER BY timestamp ASC
LIMIT @_limit::int
OFFSET @_offset::int;

-- name: QueueGetAllByStatusAndDepth :many
SELECT *
FROM queue
WHERE status = @status::int
  AND depth = @depth::int
ORDER BY timestamp ASC
LIMIT @_limit::int
OFFSET @_offset::int;

-- name: QueueGetAllByPeriod :many
SELECT *
FROM queue
WHERE timestamp >= $1 AND timestamp <= $2
ORDER BY timestamp ASC;

-- name: QueueGetAllByStatusAndUrl :one
SELECT *
FROM queue
WHERE status = @status::int
  AND url = @url::text
ORDER BY timestamp ASC LIMIT 1;

