WITH daily_posts as (
    SELECT DISTINCT
        author_id,
        TIMESTAMP_TRUNC(timestamp, DAY) as day,
        subreddit,
        count(*) as post_count
    FROM {{ ref('src_comments') }}
    GROUP BY author_id, subreddit, 2
    ORDER BY 2
),
rolling as (
    SELECT
        author_id,
        subreddit,
        day,
        SUM(post_count) OVER (
            PARTITION BY author_id, subreddit
            ORDER BY day
            ROWS BETWEEN 31 PRECEDING AND CURRENT ROW
        ) AS total_recent_post_count
    FROM daily_posts
)
SELECT * FROM rolling