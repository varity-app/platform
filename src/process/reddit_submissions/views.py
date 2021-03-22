import faust

from util.constants.kafka import Topics

from process.app import app
from .models import Submission

submissions_topic = app.topic(Topics.REDDIT_SUBMISSIONS, value_type=Submission)