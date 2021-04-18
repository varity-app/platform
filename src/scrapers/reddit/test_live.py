"""
Test the live reddit scraper
"""

from datetime import datetime, timedelta
import pytest
import asyncio
import logging

from .live import RedditScraper

logger = logging.getLogger(__name__)

START_DATE = datetime.now() - timedelta(days=7)
END_DATE = datetime.now()
LIMIT = 10

test_data = [
    ("submissions", "stocks"),
    ("comments", "stocks"),
]

@pytest.mark.asyncio
@pytest.mark.parametrize("mode,subreddit", test_data)
async def test_live_scraper(mode: str, subreddit: str) -> None:
    """Test the live reddit scraper"""

    scraper = RedditScraper(
        subreddit,
        mode,
        limit=LIMIT,
        enable_firestore=False,
        enable_publish=False,
    )

    results = await scraper.run()
    await scraper.reddit.close()
