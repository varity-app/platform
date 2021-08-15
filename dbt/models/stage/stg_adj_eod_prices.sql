SELECT
  eod_prices.symbol,
  eod_prices.date,
  eod_prices.split,
  eod_prices.close * split_factors.net_split_factor as adj_close,
  eod_prices.high * split_factors.net_split_factor as adj_high,
  eod_prices.low * split_factors.net_split_factor as adj_low,
  eod_prices.open * split_factors.net_split_factor as adj_open,
  eod_prices.volume * split_factors.net_split_factor as adj_volume,
  split_factors.net_split_factor as split_factor
FROM
  {{ ref('src_eod_prices') }} eod_prices
  LEFT JOIN (
    SELECT
      symbol,
      date,
      EXP(
        SUM(LOG(split)) OVER (
          PARTITION BY symbol
          ORDER BY
            date
        )
      ) as net_split_factor
    FROM
      {{ ref('src_eod_prices') }}
  ) split_factors
    ON eod_prices.symbol = split_factors.symbol
    AND eod_prices.date = split_factors.date