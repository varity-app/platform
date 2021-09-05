WITH comments_first_seen AS (
    SELECT
        author_id,
        subreddit,
        MIN(timestamp) as original_date
    FROM {{ ref('src_comments') }}
        GROUP BY
            author_id,
            subreddit
), submissions_first_seen AS (
    SELECT
        author_id,
        subreddit,
        MIN(timestamp) as original_date
    FROM {{ ref('src_submissions') }}
        GROUP BY
            author_id,
            subreddit
)

SELECT
    author_id,
    subreddit,
    TIMESTAMP_TRUNC(MIN(original_date), DAY) as original_date
FROM (
    SELECT * FROM comments_first_seen
    UNION ALL SELECT * FROM submissions_first_seen
) GROUP BY 1, 2