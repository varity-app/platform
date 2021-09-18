SELECT
    month,
    subreddit,
    APPROX_QUANTILES(avg_recent_comments_count, 100)[OFFSET(99)] AS percentile_99,
    APPROX_QUANTILES(avg_recent_comments_count, 100)[OFFSET(95)] AS percentile_95,
    APPROX_QUANTILES(avg_recent_comments_count, 100)[OFFSET(90)] AS percentile_90,
    APPROX_QUANTILES(avg_recent_comments_count, 100)[OFFSET(75)] AS percentile_75,
    APPROX_QUANTILES(avg_recent_comments_count, 100)[OFFSET(50)] AS percentile_50,
    APPROX_QUANTILES(avg_recent_comments_count, 100)[OFFSET(25)] AS percentile_25,
    APPROX_QUANTILES(avg_recent_comments_count, 100)[OFFSET(10)] AS percentile_10
FROM {{ ref('redditors_popularities') }}
GROUP BY month, subreddit