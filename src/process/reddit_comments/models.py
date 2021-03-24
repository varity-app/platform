"""
Faust data models for Reddit Comments
"""

import faust


class Comment(faust.Record):
    """Faust model for Reddit Comments"""

    comment_id: str
    submission_id: str
    subreddit: str
    author: str
    created_utc: str
    body: str
    upvotes: int
