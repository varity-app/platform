"""
Faust stream declarations for Kafka Topics
"""

from util.constants.kafka import Topics

from process.app import app
from .models import Comment

comments_topic = app.topic(Topics.REDDIT_COMMENTS, value_type=Comment)
