WITH percentile_10 AS (
    SELECT
        popularities.author_id,
        popularities.subreddit,
        popularities.month,
        "popularity_lowest_10_percent" AS cohort
    FROM {{ ref('redditors_popularities') }} popularities
    INNER JOIN {{ ref('redditors_popularities_quantiles') }} quantiles
        ON popularities.month = quantiles.month AND popularities.subreddit = quantiles.subreddit
    WHERE popularities.avg_recent_comments_count <= quantiles.percentile_10
), percentile_25 AS (
    SELECT
        popularities.author_id,
        popularities.subreddit,
        popularities.month,
        "popularity_lowest_25_percent" AS cohort
    FROM {{ ref('redditors_popularities') }} popularities
    INNER JOIN {{ ref('redditors_popularities_quantiles') }} quantiles
        ON popularities.month = quantiles.month AND popularities.subreddit = quantiles.subreddit
    WHERE popularities.avg_recent_comments_count <= quantiles.percentile_25
), percentile_50 AS (
    SELECT
        popularities.author_id,
        popularities.subreddit,
        popularities.month,
        "popularity_top_50_percent" AS cohort
    FROM {{ ref('redditors_popularities') }} popularities
    INNER JOIN {{ ref('redditors_popularities_quantiles') }} quantiles
        ON popularities.month = quantiles.month AND popularities.subreddit = quantiles.subreddit
    WHERE popularities.avg_recent_comments_count >= quantiles.percentile_50
), percentile_75 AS (
    SELECT
        popularities.author_id,
        popularities.subreddit,
        popularities.month,
        "popularity_top_25_percent" AS cohort
    FROM {{ ref('redditors_popularities') }} popularities
    INNER JOIN {{ ref('redditors_popularities_quantiles') }} quantiles
        ON popularities.month = quantiles.month AND popularities.subreddit = quantiles.subreddit
    WHERE popularities.avg_recent_comments_count >= quantiles.percentile_75
), percentile_90 AS (
    SELECT
        popularities.author_id,
        popularities.subreddit,
        popularities.month,
        "popularity_top_10_percent" AS cohort
    FROM {{ ref('redditors_popularities') }} popularities
    INNER JOIN {{ ref('redditors_popularities_quantiles') }} quantiles
        ON popularities.month = quantiles.month AND popularities.subreddit = quantiles.subreddit
    WHERE popularities.avg_recent_comments_count >= quantiles.percentile_90
)

SELECT * FROM percentile_10
UNION ALL SELECT * FROM percentile_25
UNION ALL SELECT * FROM percentile_50
UNION ALL SELECT * FROM percentile_75
UNION ALL SELECT * FROM percentile_90