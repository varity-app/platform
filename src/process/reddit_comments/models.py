import faust

class Comment(faust.Record):
    comment_id: str
    submission_id: str
    subreddit: str
    author: str
    created_utc: str
    body: str
    upvotes: int