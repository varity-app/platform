"""
Scrape historical reddit posts for a given day
"""

from typing import List, Tuple
import argparse
from datetime import datetime
import logging

from scrapers.reddit.historical import HistoricalRedditScraper
from util.logging import set_log_config
from util.constants.reddit import Misc
from util.constants.logging import Sources

set_log_config()
logger = logging.getLogger(Sources.HISTORICAL_REDDIT_SCRAPER)

POPULARITY_YEAR = 2020  # Year WSB got popular
WALLSTREETBETS = "wallstreetbets"


def parse_args() -> Tuple[str, List[str], int, int, List[int], int, int, int, bool]:
    """Parse arguments from command line"""
    parser = argparse.ArgumentParser()

    parser.add_argument("mode", choices=[Misc.COMMENTS, Misc.SUBMISSIONS])
    parser.add_argument("-s", "--subreddit", default="wallstreetbets")
    parser.add_argument("-y", "--year", type=int, default=2020)
    parser.add_argument("-m", "--month", type=int, default=1)
    parser.add_argument("-d", "--day", type=int, default=-1)
    parser.add_argument("-c", "--chunks", type=int, default=30)
    parser.add_argument("--limit", default=100, type=int)
    parser.add_argument("--debug", action="store_true")

    args = parser.parse_args()

    year, month, day = args.year, args.month, args.day

    subreddit = args.subreddit
    if not subreddit:
        raise ValueError(f"Invalid value for subreddit: {subreddit}")

    assert year > 2010

    if month not in range(1, 13):
        raise ValueError(f"Invalid value for month: {month}")

    if 1 <= day <= 31:
        days = [day]
    elif day == -1:
        days = range(1, 32)
    else:
        raise ValueError(f"Invalid value for day: {day}")

    assert 1 <= args.chunks <= 60

    return args.mode, subreddit, year, month, days, args.limit, args.chunks, args.debug


def main():
    """Entrypoint method"""
    mode, subreddit, year, month, days, limit, chunks, debug = parse_args()

    scraper = HistoricalRedditScraper(
        subreddit,
        mode,
        limit=limit,
        enable_firestore=not debug,
        enable_publish=not debug,
    )

    # Scrape every hour of every day
    for day in days:
        for hour in range(24):
            if year < POPULARITY_YEAR or subreddit.lower() != WALLSTREETBETS:
                start_date = datetime(
                    year=year, month=month, day=day, hour=hour, minute=0
                )
                end_date = datetime(
                    year=year, month=month, day=day, hour=hour, minute=59, second=59
                )

                logger.info(
                    f"Scraping {mode} on r/{scraper.subreddit} for"
                    f" {year}-{month}-{day} {hour}:00:00..."
                )
                num_results = scraper.run(start_date, end_date)
                logger.info(f"Saved {num_results} {mode}")

            else:  # Scrape at a much lower resolution for r/wallstreetbets
                interval = 60 // chunks
                for start_minute in range(0, 60, interval):
                    start_date = datetime(
                        year=year, month=month, day=day, hour=hour, minute=start_minute
                    )
                    end_date = datetime(
                        year=year,
                        month=month,
                        day=day,
                        hour=hour,
                        minute=start_minute + interval - 1,
                        second=59,
                    )

                    logger.info(
                        f"Scraping {mode} on r/{scraper.subreddit} for"
                        f"{year}-{month}-{day} {hour}:{start_minute}:00..."
                    )
                    num_results = scraper.run(start_date, end_date)
                    logger.info(f"Saved {num_results} {mode}")


if __name__ == "__main__":
    main()
