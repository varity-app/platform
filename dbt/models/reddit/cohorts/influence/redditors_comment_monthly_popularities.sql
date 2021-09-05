WITH comment_counts AS (
    SELECT DISTINCT
        submission_id AS parent_id,
        COUNT(*) AS counts
    FROM {{ ref('src_comments') }}
    WHERE REGEXP_CONTAINS(submission_id, r"t1_*")
    GROUP BY submission_id
), posts_with_comment_counts AS (
    SELECT
        posts.author_id,
        posts.comment_id,
        posts.subreddit,
        posts.timestamp,
        comment_counts.counts AS comment_counts
    FROM {{ ref('src_comments') }} posts
    LEFT JOIN comment_counts
        ON ("t1_" || posts.comment_id) = comment_counts.parent_id
), daily_posts AS (
    SELECT DISTINCT
        author_id,
        TIMESTAMP_TRUNC(timestamp, MONTH) AS month,
        subreddit,
        SUM(comment_counts) AS total_comment_counts
    FROM posts_with_comment_counts
    GROUP BY author_id, subreddit, 2
    ORDER BY 2
)

SELECT
    author_id,
    subreddit,
    month,
    AVG(total_comment_counts) AS avg_recent_comment_count
FROM daily_posts
GROUP BY 1, 2, 3