select
    comment_id,
    submission_id,
    subreddit,
    author_id,
    timestamp
from {{var('dataset')}}.reddit_comments_v2