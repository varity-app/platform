SELECT
    day,
    subreddit,
    APPROX_QUANTILES(age, 100)[OFFSET(90)] AS percentile_90,
    APPROX_QUANTILES(age, 100)[OFFSET(50)] AS percentile_50,
    APPROX_QUANTILES(age, 100)[OFFSET(10)] AS percentile_10
FROM {{ ref('redditors_comment_ages') }}
GROUP BY day, subreddit