WITH comments AS (
    SELECT
        parent_id,
        count(distinct symbol) = 1 AS targeted,
        (sum(question_mark_count) / sum(word_count)) > 0.1 AS inquisitive,
        APPROX_TOP_COUNT(symbol, 1)[OFFSET(0)].value AS top_symbol
    FROM {{ ref('src_comments_mentions') }}
    GROUP BY
        parent_id
)

SELECT DISTINCT
    mentions.symbol,
    mentions.timestamp,
    mentions.parent_id AS comment_id,
    comments.targeted AND mentions.symbol = comments.top_symbol AS targeted,
    comments.inquisitive
FROM {{ ref('src_comments_mentions') }} mentions
LEFT JOIN comments
    ON mentions.parent_id = comments.parent_id
    