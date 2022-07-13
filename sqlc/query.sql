-- name: GetWagers :many
SELECT
	*
FROM
	wagers
ORDER BY
	placed_at DESC
LIMIT
	? OFFSET ?;

-- name: GetWager :one
SELECT
	*
FROM
	wagers
WHERE
	id = ?
LIMIT
	1;

-- name: CreateWager :execresult
INSERT INTO
	wagers (
		total_wager_value,
		odds,
		selling_percentage,
		selling_price,
		current_selling_price,
		percentage_sold,
		amount_sold,
		placed_at
	)
VALUES
	(?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateWager :execresult
UPDATE
	wagers
SET
	current_selling_price = ?,
	percentage_sold = ?,
	amount_sold = ?
WHERE
	id = ?;

-- name: CreatePurchase :execresult
INSERT INTO
	purchases (wager_id, buying_price, bought_at)
VALUES
	(?, ?, ?);