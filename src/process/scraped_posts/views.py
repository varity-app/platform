import faust

from util.constants.kafka import Topics

from process.app import app
from .models import ScrapedPost

scraped_posts_topic = app.topic(Topics.SCRAPED_POSTS, value_type=ScrapedPost)
