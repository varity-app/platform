"""
Sentiment Analysis-related constants
"""


class Constants:
    """Kafka Message field names"""

    POSITIVE = "positive"
    NEGATIVE = "negative"
    NEUTRAL = "neutral"
    COMPOUND = "compound"
    ESTIMATOR = "estimator"


class Estimators:
    """Types of estimators used for sentiment analysis"""

    VADER = "vader"
