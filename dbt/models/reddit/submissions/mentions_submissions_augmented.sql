WITH posts_titles AS (
    SELECT
        parent_id,
        count(distinct symbol) = 1 AS targeted
    FROM {{ ref('src_submissions_mentions') }}
    WHERE parent_source = 'reddit-submission-title'
    GROUP BY
        parent_id
),
posts_titles_bodies AS (
    SELECT
        parent_id,
        count(distinct symbol) = 1 AS targeted,
        (sum(question_mark_count) / sum(word_count)) > 0.1 AS inquisitive,
        APPROX_TOP_COUNT(symbol, 1)[OFFSET(0)].value AS top_symbol
    FROM {{ ref('src_submissions_mentions') }}
    GROUP BY
        parent_id
)

SELECT DISTINCT
    mentions.symbol,
    mentions.timestamp,
    mentions.parent_id AS submission_id,
    (posts_titles.targeted OR posts_titles_bodies.targeted) AND mentions.symbol = posts_titles_bodies.top_symbol AS targeted,
    posts_titles_bodies.inquisitive
FROM {{ ref('src_submissions_mentions') }} mentions
LEFT JOIN posts_titles
    ON mentions.parent_id = posts_titles.parent_id
LEFT JOIN posts_titles_bodies
    ON mentions.parent_id = posts_titles_bodies.parent_id
    