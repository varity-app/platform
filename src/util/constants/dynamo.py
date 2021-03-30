"""
Constants related to DynamoDB
"""

import os


class Constants:
    """
    Miscellanious constants for related to DynamoDB
    """

    EXPIRATION_DATE = "expiration_date"
    DEFAULT_EXPIRATION_DAYS = 7


class Tables:
    """
    Names of tables in DynamoDB
    """

    REDDIT_SUBMISSIONS = os.environ.get(
        "DYNAMO_REDDIT_SUBMISSIONS", "reddit-submissions-dev"
    )
    REDDIT_COMMENTS = os.environ.get("DYNAMO_REDDIT_COMMENTS", "reddit-comments-dev")
