WITH months as (
    SELECT DISTINCT
        TIMESTAMP_TRUNC(timestamp, MONTH) as month,
        subreddit
    FROM {{ ref('src_submissions') }}
),
first_seen as (
    SELECT
        author_id,
        subreddit,
        MIN(TIMESTAMP_TRUNC(timestamp, MONTH)) as original_date
    FROM {{ ref('src_submissions') }}
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
FULL OUTER JOIN first_seen ON first_seen.original_date = months.month AND months.subreddit = first_seen.subreddit
