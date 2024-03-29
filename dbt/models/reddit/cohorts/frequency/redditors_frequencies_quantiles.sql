SELECT
    month,
    subreddit,
    APPROX_QUANTILES(posts_count, 100)[OFFSET(90)] AS percentile_90,
    APPROX_QUANTILES(posts_count, 100)[OFFSET(75)] AS percentile_75,
    APPROX_QUANTILES(posts_count, 100)[OFFSET(50)] AS percentile_50,
    APPROX_QUANTILES(posts_count, 100)[OFFSET(25)] AS percentile_25,
    APPROX_QUANTILES(posts_count, 100)[OFFSET(10)] AS percentile_10
FROM {{ ref('redditors_frequencies') }}
GROUP BY month, subreddit