{{ config(
    materialized='view',
    partition_by={
      "field": "date",
      "data_type": "timestamp",
      "granularity": "day"
    }
)}}

WITH source_data AS (
    SELECT
        DISTINCT symbol,
        date,
        close,
        high,
        low,
        open,
        volume,
        split
    FROM
        {{var('dataset')}}.eod_prices
    
    {% if is_incremental() %} -- this filter will only be applied on an incremental run
    where
        date >= (
            SELECT
                MAX(date)
            FROM
                {{this}}
        )

    {% endif %}
)

SELECT
    *
FROM
    source_data