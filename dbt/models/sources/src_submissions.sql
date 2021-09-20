select
    submission_id,
    subreddit,
    timestamp,
    author_id
from {{var('dataset')}}.reddit_submissions_v2