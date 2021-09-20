SELECT
    *
FROM {{var('dataset')}}.ticker_mentions_v2
WHERE parent_source = 'reddit-comment'