WITH date_sequences AS (
    SELECT
        author_id,
        subreddit,
        original_date,
        GENERATE_TIMESTAMP_ARRAY(original_date, CURRENT_TIMESTAMP(), INTERVAL 30 DAY) AS months
    FROM {{ ref('redditors_first_seen') }}
)

SELECT
    author_id,
    subreddit,
    TIMESTAMP_DIFF(month, original_date, DAY) as age,
    month
FROM date_sequences
CROSS JOIN UNNEST(date_sequences.months) AS month