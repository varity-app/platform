"""
Kafka-related constants
"""

import os


class Config:
    """General Kafka configuration constants"""

    BOOTSTRAP_SERVERS = os.environ.get("BOOTSTRAP_SERVERS")
    SASL_MECHANISM = 'PLAIN'
    SECURITY_PROTOCOL = 'SASL_SSL'
    SASL_USERNAME = os.environ.get("SASL_USERNAME")
    SASL_PASSWORD = os.environ.get("SASL_PASSWORD")

    OBJ = {
        'bootstrap.servers': BOOTSTRAP_SERVERS,
        'sasl.mechanism': SASL_MECHANISM,
        'security.protocol': SECURITY_PROTOCOL,
        'sasl.username': SASL_USERNAME,
        'sasl.password': SASL_PASSWORD,
    }


class Topics:
    """Kafka topic names"""

    REDDIT_SUBMISSIONS = 'reddit-submissions'
    REDDIT_COMMENTS = 'reddit-comments'
    TICKER_MENTIONS = 'ticker-mentions'
    SCRAPED_POSTS = 'scraped-posts'
    POST_SENTIMENT = 'post-sentiment'

    LOGS = 'logs'


class Groups:
    """Kafka consumer groups"""

    SUBMISSION_CONSUMERS = 'submission-consumers'
    COMMENT_CONSUMERS = 'comment-consumers'
    SENTIMENT_ESTIMATORS = 'sentiment-estimators'
