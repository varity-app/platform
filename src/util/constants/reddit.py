"""
Reddit scraping-related constants
"""

import os


class Config:
    """Reddit API configuration variables"""

    REDDIT_USERNAME = os.environ.get("REDDIT_USERNAME")
    REDDIT_PASSWORD = os.environ.get("REDDIT_PASSWORD")
    REDDIT_CLIENT_ID = os.environ.get("REDDIT_CLIENT_ID")
    REDDIT_CLIENT_SECRET = os.environ.get("REDDIT_CLIENT_SECRET")
    REDDIT_USER_AGENT = os.environ.get("REDDIT_USER_AGENT")


class Misc:
    """Miscellanious constants"""

    SUBMISSIONS = "submissions"
    COMMENTS = "comments"


class SubmissionConstants:
    """Reddit submission object fields"""

    ID = "submission_id"
    TITLE = "title"
    CREATED_UTC = "created_utc"
    SELFTEXT = "selftext"
    UPVOTES = "upvotes"
    SUBREDDIT = "subreddit"
    AUTHOR = "author"
    IS_ORIGINAL_CONTENT = "is_original_content"
    IS_TEXT = "is_text"
    NAME = "name"
    NUM_COMMENTS = "num_comments"
    NSFW = "nsfw"
    PERMALINK = "permalink"
    URL = "url"
    UPVOTE_RATIO = "upvote_ratio"


class CommentConstants:
    """Reddit comment object fields"""

    ID = "comment_id"
    SUBMISSION_ID = "submission_id"
    BODY = "body"
    CREATED_UTC = "created_utc"
    UPVOTES = "upvotes"
    SUBREDDIT = "subreddit"
    AUTHOR = "author"
