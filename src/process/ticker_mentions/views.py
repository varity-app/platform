"""
Faust stream declarations for Kafka Topics
"""

from util.constants.kafka import Topics

from process.app import app
from .models import TickerMention

ticker_mentions_topic = app.topic(Topics.TICKER_MENTIONS, value_type=TickerMention)
