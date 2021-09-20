WITH popularities AS (
    SELECT * FROM {{ ref('redditors_comment_monthly_popularities') }}
    UNION ALL SELECT * FROM {{ ref('redditors_submission_monthly_popularities') }}
)

SELECT
    month,
    subreddit,
    author_id,
    AVG(avg_recent_comment_count) AS avg_recent_comments_count
FROM popularities
GROUP BY 1, 2, 3