"""
Faust stream declarations for Kafka Topics
"""

from util.constants.kafka import Topics

from process.app import app
from .models import SentimentEstimate

sentiment_topic = app.topic(Topics.POST_SENTIMENT, value_type=SentimentEstimate)
