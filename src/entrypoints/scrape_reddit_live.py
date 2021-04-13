"""
Scrape reddit at a given interval for new submissions or comments
"""

from time import sleep
import asyncio
import argparse
from datetime import datetime
import gc
import traceback
import logging
from asyncprawcore.exceptions import RequestException, ResponseException


from scrapers.reddit import RedditScraper
from util.logging import set_log_config
from util.constants.reddit import Misc
from util.constants.logging import Sources

set_log_config()
logger = logging.getLogger(Sources.REDDIT_SCRAPER)


def parse_args():
    """Parse arguments from command line"""
    parser = argparse.ArgumentParser()

    parser.add_argument("mode", choices=[Misc.COMMENTS, Misc.SUBMISSIONS])
    parser.add_argument("-s", "--subreddits", default="wallstreetbets")
    parser.add_argument("--sleep", default=60, type=int)
    parser.add_argument("--limit", default=200, type=int)

    args = parser.parse_args()

    subreddits = args.subreddits

    if not subreddits:
        raise ValueError(f"Invalid value for subreddits: {subreddits}")

    subreddits = subreddits.split(",")

    return args.mode, subreddits, args.sleep, args.limit


async def main():
    """Main method"""

    try:
        mode, subreddits, sleep_interval, limit = parse_args()

        scrapers = [
            RedditScraper(subreddit, mode, limit=limit) for subreddit in subreddits
        ]

        while True:
            for subreddit, scraper in zip(subreddits, scrapers):
                logger.info(
                    f"Scraping for new {mode} on r/{subreddit} at {datetime.now()}...",
                )

                try:
                    num_results = await scraper.run()
                except (RequestException, ResponseException):
                    logger.error(
                        f"Error scraping on r/{subreddit}. Will retry: {traceback.format_exc()}"
                    )
                    continue

                logger.info(f"Saved {num_results} {mode}")

            gc.collect()
            logger.info(
                f"Sleeping for {sleep_interval} seconds...",
            )
            sleep(sleep_interval)

    except Exception as ex:
        logger.error(traceback.format_exc())
        raise ex


if __name__ == "__main__":
    asyncio.run(main())
