WITH oldest_redditors AS (
    SELECT
        ages.author_id,
        ages.subreddit,
        ages.day,
        ages.age >= quantiles.percentile_90 in_percentile_90,
        ages.age >= quantiles.percentile_50 in_percentile_50,
        ages.age <= quantiles.percentile_10 in_lowest_percentile_10
    FROM {{ ref('redditors_comment_ages') }} ages
    INNER JOIN {{ ref('redditors_comment_ages_quantiles') }} quantiles
        ON ages.day = quantiles.day AND ages.subreddit = quantiles.subreddit
),
mentions_with_authors AS (
    SELECT
        mentions.symbol,
        mentions.timestamp,
        TIMESTAMP_TRUNC(mentions.timestamp, DAY) as day,
        mentions.targeted,
        mentions.inquisitive,
        comments.author_id,
        comments.subreddit
    FROM {{ ref('mentions_comments_augmented') }} mentions
    LEFT JOIN {{ ref('src_comments') }} comments
        ON mentions.comment_id = comments.comment_id
)

SELECT
    mentions_with_authors.symbol,
    mentions_with_authors.subreddit,
    TIMESTAMP_TRUNC(mentions_with_authors.timestamp, HOUR) AS hour,
    COUNTIF(oldest_redditors.in_percentile_90) AS count_in_percentile_90,
    COUNTIF(mentions_with_authors.targeted AND oldest_redditors.in_percentile_90) AS count_targeted_in_percentile_90,
    COUNTIF(mentions_with_authors.inquisitive AND oldest_redditors.in_percentile_90) AS count_inquisitive_in_percentile_90,
    COUNTIF(oldest_redditors.in_percentile_50) AS count_in_percentile_50,
    COUNTIF(mentions_with_authors.targeted AND oldest_redditors.in_percentile_50) AS count_targeted_in_percentile_50,
    COUNTIF(mentions_with_authors.inquisitive AND oldest_redditors.in_percentile_50) AS count_inquisitive_in_percentile_50,
    COUNTIF(oldest_redditors.in_lowest_percentile_10) AS count_in_lowest_percentile_10,
    COUNTIF(mentions_with_authors.targeted AND oldest_redditors.in_lowest_percentile_10) AS count_targeted_in_lowest_percentile_10,
    COUNTIF(mentions_with_authors.inquisitive AND oldest_redditors.in_lowest_percentile_10) AS count_inquisitive_in_lowest_percentile_10
FROM mentions_with_authors
LEFT JOIN oldest_redditors
    ON mentions_with_authors.day = oldest_redditors.day
        AND mentions_with_authors.subreddit = oldest_redditors.subreddit
        AND mentions_with_authors.author_id = oldest_redditors.author_id
GROUP BY 1, 2, 3
ORDER BY 3 DESC