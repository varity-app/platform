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

    REDDIT_SUBMISSIONS = f"reddit-submissions-beam-{DEPLOYMENT}"
    REDDIT_COMMENTS = f"reddit-comments-{DEPLOYMENT}"
    TICKER_MENTIONS = f"ticker-mentions-{DEPLOYMENT}"
    SCRAPED_POSTS = f"scraped-posts-{DEPLOYMENT}"


class BeamSubscriptions:
    """Pub/Sub subscriptions created for beam pipelines"""

    PREFIX = f"projects/{Config.PROJECT}/subscriptions"

    REDDIT_SUBMISSIONS = f"reddit-submissions-beam-{DEPLOYMENT}"
    REDDIT_COMMENTS = f"reddit-comments-beam-{DEPLOYMENT}"
    TICKER_MENTIONS = f"ticker-mentions-beam-{DEPLOYMENT}"
    SCRAPED_POSTS = f"scraped-posts-beam-{DEPLOYMENT}"
