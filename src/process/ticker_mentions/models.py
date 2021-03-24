"""
Faust data models for Ticker Mentions
"""

import faust


class TickerMention(faust.Record):
    """Faust model for Ticker Mentions"""

    stock_name: str
    data_source: str
    parent_source: str
    parent_id: str
    created_utc: str
    mention_type: str
