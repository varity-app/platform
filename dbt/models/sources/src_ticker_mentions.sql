with source_data as (
    select * from {{var('dataset')}}.ticker_mentions_v2
)

select *
from source_data