WITH comment_mentions AS (
    SELECT
        mentions.symbol,
        mentions.timestamp,
        mentions.targeted,
        mentions.inquisitive,
        comments.author_id,
        comments.subreddit,
        "reddit-comment" AS source
    FROM {{ ref('mentions_comments_augmented') }} mentions
    LEFT JOIN {{ ref('src_comments') }} comments
        ON mentions.comment_id = comments.comment_id
), submission_mentions AS (
    SELECT
        mentions.symbol,
        mentions.timestamp,
        mentions.targeted,
        mentions.inquisitive,
        submissions.author_id,
        submissions.subreddit,
        "reddit-submission" AS source
    FROM {{ ref('mentions_submissions_augmented') }} mentions
    LEFT JOIN {{ ref('src_submissions') }} submissions
        ON mentions.submission_id = submissions.submission_id
)

SELECT * FROM comment_mentions
UNION ALL SELECT * FROM submission_mentions
