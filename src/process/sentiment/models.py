"""
Faust data models for Sentiment Estimates
"""

import faust


class SentimentEstimate(faust.Record):
    """Faust model for a Sentiment Estimate"""

    data_source: str
    parent_source: str
    parent_id: str
    estimator: str
    compound: float
    positive: float
    neutral: float
    negative: float
