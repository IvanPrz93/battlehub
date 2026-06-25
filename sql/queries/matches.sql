-- name: CreateMatch :one
INSERT INTO matches (id, player1_id, player2_id, created_at)
VALUES (
    gen_random_uuid ( ),
    $1,
    $2,
    NOW()
)
RETURNING *;

-- name: UpdateMatchResults :one
UPDATE matches
SET winner_id = $2, completed_at = NOW()
WHERE id = $1 AND winner_id IS NULL AND completed_at IS NULL
RETURNING *;

-- name: GetMatchHistory :many
SELECT * from matches
WHERE player1_id = $1 OR player2_id = $1
ORDER BY created_at DESC;
