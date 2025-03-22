-- name: SaveCandle :exec
INSERT INTO candles (
  symbol,
  open_time,
  open,
  high,
  low,
  close,
  volume,
  close_time,
  quote_volume,
  count,
  taker_buy_volume,
  taker_buy_quote_volume
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetLatestCandle :one
SELECT *
FROM candles
WHERE symbol = ?
ORDER BY open_time DESC
LIMIT 1;

-- name: GetCandles :many
SELECT * FROM (
  SELECT * 
  FROM candles
  WHERE symbol = ?
  ORDER BY open_time DESC
  LIMIT ?
) AS c
ORDER BY open_time ASC