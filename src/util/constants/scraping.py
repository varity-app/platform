"""
General scraping-related constants
"""


class ScrapedPostConstants:
    """Fields for a ScrapedPost Kafka message"""

    DATA_SOURCE = "data_source"
    PARENT_SOURCE = "parent_source"
    PARENT_ID = "parent_id"
    TEXT = "text"
    TIMESTAMP = "timestamp"


class TickerMentionsConstants:
    """Fields for a TickerMention Kafka message"""

    STOCK_NAME = "stock_name"
    PARENT_SOURCE = "parent_source"
    PARENT_ID = "parent_id"
    DATA_SOURCE = "data_source"
    CREATED_UTC = "created_utc"
    MENTION_TYPE = "mention_type"


class MentionTypes:
    """All accepted values for the `mention-type` field in a TickerMention Kafka message"""

    TICKER = "ticker"


class DataSources:
    """Valid values for the `data_source` field"""

    REDDIT = "reddit"


class ParentSources:
    """Valid values for the `parent_source` field"""

    COMMENT_BODY = "comment_body"
    SUBMISSION_TITLE = "submission_title"
    SUBMISSION_SELFTEXT = "submission_selftext"
