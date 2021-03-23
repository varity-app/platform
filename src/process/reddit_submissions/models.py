"""
Faust data models for Reddit Submissions
"""

import faust


class Submission(faust.Record):
    """Faust model for Reddit Submissions"""

    submission_id: str
    subreddit: str
    title: str
    created_utc: str
    name: str
    selftext: str
    author: str
    is_original_content: bool
    is_text: bool
    nsfw: bool
    num_comments: int
    permalink: str
    upvotes: int
    upvote_ratio: float
    url: str
