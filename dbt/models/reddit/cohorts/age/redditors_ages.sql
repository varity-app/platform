WITH monthly_submissions AS (
    SELECT DISTINCT
        author_id,
        subreddit,
        TIMESTAMP_TRUNC(timestamp, MONTH) as month,
    FROM {{ ref('src_submissions') }}
), monthly_comments AS (
    SELECT DISTINCT
        author_id,
        subreddit,
        TIMESTAMP_TRUNC(timestamp, MONTH) as month,
    FROM {{ ref('src_comments') }}
), monthly_posts AS (
    SELECT DISTINCT
        *
    FROM (
        SELECT * FROM monthly_submissions
        UNION ALL SELECT * FROM monthly_comments
    )
)

SELECT
    monthly_posts.author_id,
    monthly_posts.subreddit,
    TIMESTAMP_DIFF(monthly_posts.month, first_seen.original_date, DAY) as age,
    monthly_posts.month
FROM monthly_posts
LEFT JOIN {{ ref('redditors_first_seen') }} as first_seen
    ON monthly_posts.author_id = first_seen.author_id
        AND monthly_posts.subreddit = first_seen.subreddit
