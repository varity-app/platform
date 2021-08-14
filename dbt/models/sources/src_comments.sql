with source_data as (
    select * from {{var('dataset')}}.reddit_comments_v2
)

select *
from source_data