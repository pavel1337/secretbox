-- name: getSecret :one
SELECT * FROM secrets WHERE id = ? LIMIT 1;

-- name: createSecret :execresult
INSERT INTO secrets (
  id, content, passphrase, expires_at
) VALUES (
  ?, ?, ?, ?
);

-- name: deleteSecret :exec
DELETE FROM secrets
WHERE id = ?;

-- name: listExpiredSecretsIds :many
SELECT id FROM secrets
WHERE expires_at < ?
ORDER BY id;
