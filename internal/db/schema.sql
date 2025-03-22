CREATE TABLE `candles` (
  `symbol` varchar(20) NOT NULL,
  `open_time` bigint NOT NULL,
  `open` decimal(18,8) NOT NULL,
  `high` decimal(18,8) NOT NULL,
  `low` decimal(18,8) NOT NULL,
  `close` decimal(18,8) NOT NULL,
  `volume` decimal(18,8) NOT NULL,
  `close_time` bigint NOT NULL,
  `quote_volume` decimal(18,8) NOT NULL,
  `count` bigint NOT NULL,
  `taker_buy_volume` decimal(18,8) NOT NULL,
  `taker_buy_quote_volume` decimal(18,8) NOT NULL,
  PRIMARY KEY (`symbol`,`open_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;