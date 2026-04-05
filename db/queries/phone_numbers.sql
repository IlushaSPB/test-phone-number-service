-- name: InsertPhoneNumber :execrows
INSERT INTO phone_numbers (phone_number, source, country, region, provider)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (phone_number) DO NOTHING;

-- name: SearchPhoneNumbers :many
SELECT id, phone_number, source, country, region, provider, created_at
FROM phone_numbers
WHERE
    (sqlc.narg('number')::TEXT IS NULL OR phone_number LIKE '%' || sqlc.narg('number') || '%')
    AND (sqlc.narg('country')::TEXT IS NULL OR country = sqlc.narg('country'))
    AND (sqlc.narg('region')::TEXT IS NULL OR region = sqlc.narg('region'))
    AND (sqlc.narg('provider')::TEXT IS NULL OR provider = sqlc.narg('provider'))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountPhoneNumbers :one
SELECT COUNT(*)
FROM phone_numbers
WHERE
    (sqlc.narg('number')::TEXT IS NULL OR phone_number LIKE '%' || sqlc.narg('number') || '%')
    AND (sqlc.narg('country')::TEXT IS NULL OR country = sqlc.narg('country'))
    AND (sqlc.narg('region')::TEXT IS NULL OR region = sqlc.narg('region'))
    AND (sqlc.narg('provider')::TEXT IS NULL OR provider = sqlc.narg('provider'));
