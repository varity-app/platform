"""
Constants related to BigQuery
"""

from . import DEPLOYMENT


class BatchSizes:
    """Size of batches of streaming inserts to BigQuery"""

    REDDIT_SUBMISSIONS = 25
    REDDIT_COMMENTS = 100
    SCRAPED_POSTS = 100


class Tables:
    """IDs of tables inside BigQuery"""

    DB = f"varity:scraping_{DEPLOYMENT}"

    REDDIT_SUBMISSIONS = f"{DB}.reddit_submissions"
    REDDIT_COMMENTS = f"{DB}.reddit_comments"
    SCRAPED_POSTS = f"{DB}.scraped_posts"
    TICKER_MENTIONS = f"{DB}.ticker_mentions"
