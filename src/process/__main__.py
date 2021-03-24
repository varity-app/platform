"""
Entrypoint file for Faust
"""

# pylint: disable=W0611
from .app import app
from . import reddit_submissions
from . import reddit_comments
from . import scraped_posts

app.main()
