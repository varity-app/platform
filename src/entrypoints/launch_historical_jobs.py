"""
Helper script for deploying historical reddit jobs on kubernetes
"""

from typing import Tuple
import argparse
import logging

from k8s import init_k8s, create_job
from k8s.reddit_historical import create_scraper_job_object

from util.constants.reddit import Misc
from util.logging import set_log_config

set_log_config()
logger = logging.getLogger(__name__)


def parse_args() -> Tuple[str, str, int, int, int]:
    """Parse arguments from command line"""
    parser = argparse.ArgumentParser()

    parser.add_argument("mode", choices=[Misc.COMMENTS, Misc.SUBMISSIONS])
    parser.add_argument("-s", "--subreddit", default="wallstreetbets")
    parser.add_argument("-y", "--year", type=int, default=2020)
    parser.add_argument("-m", "--month", type=int, default=1)
    parser.add_argument("-d", "--day", type=int, default=-1)

    args = parser.parse_args()

    year, month, day = args.year, args.month, args.day

    subreddit = args.subreddit
    if not subreddit:
        raise ValueError(f"Invalid value for subreddit: {subreddit}")

    assert year > 2010

    if month not in range(1, 13):
        raise ValueError(f"Invalid value for month: {month}")

    assert day == -1 or 1 <= day <= 31

    return args.mode, subreddit, year, month, day


def main():
    """Entrypoint method"""

    # Load args
    mode, subreddit, year, month, day = parse_args()

    # Initialize k8s
    api_instance = init_k8s()

    # Create job spec
    body = create_scraper_job_object(
        mode, subreddit, scraping_year=year, scraping_month=month, scraping_day=day
    )

    # Create job
    create_job(api_instance, body)

    if day == -1:
        logger.info(
            f"Created scraping job for {mode} on r/{subreddit} for {year}-{month}"
        )
    else:
        logger.info(
            f"Created scraping job for {mode} on r/{subreddit} for {year}-{month}-{day}"
        )


if __name__ == "__main__":
    main()
