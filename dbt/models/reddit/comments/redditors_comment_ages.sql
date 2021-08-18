WITH daily_posts as (
    SELECT DISTINCT
        author_id,
        TIMESTAMP_TRUNC(timestamp, DAY) as day,
        subreddit
    FROM {{ ref('src_comments') }}
),
first_seen as (
    SELECT
        author_id,
        subreddit,
        MIN(timestamp) as original_date
    FROM {{ ref('src_comments') }}
        GROUP BY
            author_id,
            subreddit
)

SELECT
    daily_posts.author_id,
    daily_posts.day,
    daily_posts.subreddit,
    TIMESTAMP_DIFF(daily_posts.day, first_seen.original_date, DAY) as age
FROM
    daily_posts
LEFT JOIN first_seen
    ON daily_posts.author_id = first_seen.author_id
        AND daily_posts.subreddit = first_seen.subreddit

