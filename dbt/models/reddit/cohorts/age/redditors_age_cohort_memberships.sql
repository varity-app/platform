WITH percentile_10 AS (
    SELECT
        ages.author_id,
        ages.subreddit,
        ages.month,
        "age_lowest_10_percent" AS cohort
    FROM {{ ref('redditors_ages') }} ages
    INNER JOIN {{ ref('redditors_ages_quantiles') }} quantiles
        ON ages.month = quantiles.month AND ages.subreddit = quantiles.subreddit
    WHERE ages.age <= quantiles.percentile_10
), percentile_25 AS (
    SELECT
        ages.author_id,
        ages.subreddit,
        ages.month,
        "age_lowest_25_percent" AS cohort
    FROM {{ ref('redditors_ages') }} ages
    INNER JOIN {{ ref('redditors_ages_quantiles') }} quantiles
        ON ages.month = quantiles.month AND ages.subreddit = quantiles.subreddit
    WHERE ages.age <= quantiles.percentile_25
), percentile_50 AS (
    SELECT
        ages.author_id,
        ages.subreddit,
        ages.month,
        "age_top_50_percent" AS cohort
    FROM {{ ref('redditors_ages') }} ages
    INNER JOIN {{ ref('redditors_ages_quantiles') }} quantiles
        ON ages.month = quantiles.month AND ages.subreddit = quantiles.subreddit
    WHERE ages.age >= quantiles.percentile_50
), percentile_75 AS (
    SELECT
        ages.author_id,
        ages.subreddit,
        ages.month,
        "age_top_25_percent" AS cohort
    FROM {{ ref('redditors_ages') }} ages
    INNER JOIN {{ ref('redditors_ages_quantiles') }} quantiles
        ON ages.month = quantiles.month AND ages.subreddit = quantiles.subreddit
    WHERE ages.age >= quantiles.percentile_75
), percentile_90 AS (
    SELECT
        ages.author_id,
        ages.subreddit,
        ages.month,
        "age_top_10_percent" AS cohort
    FROM {{ ref('redditors_ages') }} ages
    INNER JOIN {{ ref('redditors_ages_quantiles') }} quantiles
        ON ages.month = quantiles.month AND ages.subreddit = quantiles.subreddit
    WHERE ages.age >= quantiles.percentile_90
)

SELECT * FROM percentile_10
UNION ALL SELECT * FROM percentile_25
UNION ALL SELECT * FROM percentile_50
UNION ALL SELECT * FROM percentile_75
UNION ALL SELECT * FROM percentile_90

