import faust

class SentimentEstimate(faust.Record):
    data_source: str
    parent_source: str
    parent_id: str
    estimator: str
    compound: float
    positive: float
    neutral: float
    negative: float