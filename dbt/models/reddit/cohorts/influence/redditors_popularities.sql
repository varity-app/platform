SELECT
    comments.day,
    comments.subreddit,
    comments.author_id,
    comments.total_recent_comment_count + submissions.total_recent_comment_count AS comments_count
FROM {{ ref('redditors_comment_monthly_popularities') }} comments
LEFT JOIN {{ ref('redditors_submission_monthly_popularities') }} submissions
    ON comments.day = submissions.day
        AND comments.subreddit = submissions.subreddit
        AND comments.author_id = submissions.author_id