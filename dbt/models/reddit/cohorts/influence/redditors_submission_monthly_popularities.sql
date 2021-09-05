WITH comment_counts as (
    SELECT DISTINCT
        submission_id,
        count(*) as counts
    FROM {{ ref('src_comments') }}
    WHERE REGEXP_CONTAINS(submission_id, r"t3_*")
    GROUP BY submission_id
),
posts_with_comment_counts as (
    SELECT
        posts.author_id,
        posts.submission_id,
        posts.subreddit,
        posts.timestamp,
        comment_counts.counts as comment_counts
    FROM {{ ref('src_submissions') }} posts
    LEFT JOIN comment_counts
        ON ("t3_" || posts.submission_id) = comment_counts.submission_id
),
daily_posts as (
SELECT DISTINCT
        author_id,
        TIMESTAMP_TRUNC(timestamp, DAY) as day,
        subreddit,
        sum(comment_counts) as total_comment_counts
    FROM posts_with_comment_counts
    GROUP BY author_id, subreddit, 2
    ORDER BY 2
)

SELECT
    author_id,
    subreddit,
    day,
    SUM(total_comment_counts) OVER (
        PARTITION BY author_id, subreddit
        ORDER BY day
        ROWS BETWEEN 31 PRECEDING AND CURRENT ROW
    ) AS total_recent_comment_count
FROM daily_posts