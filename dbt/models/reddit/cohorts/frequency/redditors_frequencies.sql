WITH comment_frequencies AS (
    SELECT
        author_id,
        subreddit,
        TIMESTAMP_TRUNC(timestamp, MONTH) as month,
        count(*) as posts_count
    FROM {{ ref('src_comments') }}
    GROUP BY author_id, subreddit, 3
), submission_frequencies AS (
    SELECT
        author_id,
        subreddit,
        TIMESTAMP_TRUNC(timestamp, MONTH) as month,
        count(*) as posts_count
    FROM {{ ref('src_submissions') }}
    GROUP BY author_id, subreddit, 3
), all_posts AS (
    SELECT * FROM comment_frequencies
    UNION ALL SELECT * FROM submission_frequencies
)

SELECT
    author_id,
    subreddit,
    month,
    SUM(posts_count) AS posts_count
FROM all_posts
GROUP BY 1, 2, 3