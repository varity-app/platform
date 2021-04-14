"""
Pub/Sub-related constants
"""

from . import DEPLOYMENT


class Config:
    """General Pub/Sub configuration constants"""

    PROJECT = "varity"


class Topics:
    """Pub/Sub topic names"""

    PREFIX = f"projects/{Config.PROJECT}/topics"

    REDDIT_SUBMISSIONS = f"reddit-submissions-{DEPLOYMENT}"
    REDDIT_COMMENTS = f"reddit-comments-{DEPLOYMENT}"
    TICKER_MENTIONS = f"ticker-mentions-{DEPLOYMENT}"
    SCRAPED_POSTS = f"scraped-posts-{DEPLOYMENT}"
    POST_SENTIMENT = f"post-sentiment-{DEPLOYMENT}"

    LOGS = "logs"


class Groups:
    """Pub/Sub consumer groups"""

    SUBMISSION_CONSUMERS = "submission-consumers"
    COMMENT_CONSUMERS = "comment-consumers"
    SENTIMENT_ESTIMATORS = "sentiment-estimators"
