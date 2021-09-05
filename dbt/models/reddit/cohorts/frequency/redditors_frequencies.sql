SELECT
    comments.day,
    comments.subreddit,
    comments.author_id,
    comments.total_recent_post_count + submissions.total_recent_post_count AS posts_count
FROM {{ ref('redditors_comment_monthly_frequencies') }} comments
LEFT JOIN {{ ref('redditors_submission_monthly_frequencies') }} submissions
    ON comments.day = submissions.day
        AND comments.subreddit = submissions.subreddit
        AND comments.author_id = submissions.author_id