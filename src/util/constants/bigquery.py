"""
Constants related to BigQuery
"""

from . import DEPLOYMENT


class Tables:
    """IDs of tables inside BigQuery"""

    DB = f"varity:scraping_{DEPLOYMENT}"

    REDDIT_SUBMISSIONS = f"{DB}.reddit_submissions"
    REDDIT_COMMENTS = f"{DB}.reddit_comments"
    SCRAPED_POSTS = f"{DB}.scraped_posts"
    TICKER_MENTIONS = f"{DB}.ticker_mentions"
