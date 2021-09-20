WITH percentile_10 AS (
    SELECT
        frequencies.author_id,
        frequencies.subreddit,
        frequencies.month,
        "post_frequency_lowest_10_percent" AS cohort
    FROM {{ ref('redditors_frequencies') }} frequencies
    INNER JOIN {{ ref('redditors_frequencies_quantiles') }} quantiles
        ON frequencies.month = quantiles.month AND frequencies.subreddit = quantiles.subreddit
    WHERE frequencies.posts_count <= quantiles.percentile_10
), percentile_25 AS (
    SELECT
        frequencies.author_id,
        frequencies.subreddit,
        frequencies.month,
        "post_frequency_lowest_25_percent" AS cohort
    FROM {{ ref('redditors_frequencies') }} frequencies
    INNER JOIN {{ ref('redditors_frequencies_quantiles') }} quantiles
        ON frequencies.month = quantiles.month AND frequencies.subreddit = quantiles.subreddit
    WHERE frequencies.posts_count <= quantiles.percentile_25
), percentile_50 AS (
    SELECT
        frequencies.author_id,
        frequencies.subreddit,
        frequencies.month,
        "post_frequency_top_50_percent" AS cohort
    FROM {{ ref('redditors_frequencies') }} frequencies
    INNER JOIN {{ ref('redditors_frequencies_quantiles') }} quantiles
        ON frequencies.month = quantiles.month AND frequencies.subreddit = quantiles.subreddit
    WHERE frequencies.posts_count >= quantiles.percentile_50
), percentile_75 AS (
    SELECT
        frequencies.author_id,
        frequencies.subreddit,
        frequencies.month,
        "post_frequency_top_25_percent" AS cohort
    FROM {{ ref('redditors_frequencies') }} frequencies
    INNER JOIN {{ ref('redditors_frequencies_quantiles') }} quantiles
        ON frequencies.month = quantiles.month AND frequencies.subreddit = quantiles.subreddit
    WHERE frequencies.posts_count >= quantiles.percentile_75
), percentile_90 AS (
    SELECT
        frequencies.author_id,
        frequencies.subreddit,
        frequencies.month,
        "post_frequency_top_10_percent" AS cohort
    FROM {{ ref('redditors_frequencies') }} frequencies
    INNER JOIN {{ ref('redditors_frequencies_quantiles') }} quantiles
        ON frequencies.month = quantiles.month AND frequencies.subreddit = quantiles.subreddit
    WHERE frequencies.posts_count >= quantiles.percentile_90
)

SELECT * FROM percentile_10
UNION ALL SELECT * FROM percentile_25
UNION ALL SELECT * FROM percentile_50
UNION ALL SELECT * FROM percentile_75
UNION ALL SELECT * FROM percentile_90