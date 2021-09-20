WITH months as (
    SELECT DISTINCT
        TIMESTAMP_TRUNC(timestamp, DAY) as month,
        subreddit
    FROM {{ ref('src_comments') }}
),
first_seen as (
    SELECT
        author_id,
        subreddit,
        MIN(TIMESTAMP_TRUNC(timestamp, DAY)) as original_date
    FROM {{ ref('src_comments') }}
        GROUP BY
            author_id,
            subreddit
)

SELECT
    first_seen.author_id,
    months.month,
    first_seen.subreddit,
    TIMESTAMP_DIFF(months.month, first_seen.original_date, DAY) as age
FROM
    months
FULL OUTER JOIN first_seen ON months.subreddit = first_seen.subreddit
