"""
Test the historical reddit scraper
"""

from datetime import datetime, timedelta
import pytest

from .historical import HistoricalRedditScraper


START_DATE = datetime.now() - timedelta(days=7)
END_DATE = datetime.now()
LIMIT = 10

test_data = [
    ("submissions", "stocks"),
    ("comments", "stocks"),
]


@pytest.mark.parametrize("mode,subreddit", test_data)
def test_historical_scraper(mode: str, subreddit: str) -> None:
    """Test the historical reddit scraper"""

    scraper = HistoricalRedditScraper(
        subreddit,
        mode,
        limit=LIMIT,
        enable_firestore=False,
        enable_publish=False,
    )

    scraper.run(START_DATE, END_DATE)
