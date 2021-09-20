WITH comment_mentions AS (
    SELECT
        mentions.symbol,
        mentions.timestamp,
        CASE WHEN mentions.targeted IS NULL THEN FALSE ELSE mentions.targeted END targeted,
        CASE WHEN mentions.inquisitive IS NULL THEN FALSE ELSE mentions.inquisitive END inquisitive,
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
        CASE WHEN mentions.targeted IS NULL THEN FALSE ELSE mentions.targeted END targeted,
        CASE WHEN mentions.inquisitive IS NULL THEN FALSE ELSE mentions.inquisitive END inquisitive,
        submissions.author_id,
        submissions.subreddit,
        "reddit-submission" AS source
    FROM {{ ref('mentions_submissions_augmented') }} mentions
    LEFT JOIN {{ ref('src_submissions') }} submissions
        ON mentions.submission_id = submissions.submission_id
)

SELECT * FROM comment_mentions
UNION ALL SELECT * FROM submission_mentions
