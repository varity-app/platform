"Scrapes reddit at a given interval for new submissions or comments"

from time import sleep
import asyncio
import argparse
from datetime import datetime
import sys
import gc
from asyncprawcore.exceptions import RequestException
import logging

from scrapers.reddit import RedditScraper
from util.constants.reddit import Misc

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def parse_args():
    '''Parse arguments from command line'''
    parser = argparse.ArgumentParser()
    
    parser.add_argument('mode', choices=[Misc.COMMENTS, Misc.SUBMISSIONS])
    parser.add_argument('-s', '--subreddits', default='wallstreetbets')
    parser.add_argument('--sleep', default=60, type=int)
    parser.add_argument('--limit', default=200, type=int)

    args = parser.parse_args()

    subreddits = args.subreddits

    if not subreddits:
        raise ValueError(f'Invalid value for subreddits: {subreddits}')

    subreddits = subreddits.split(',')

    return args.mode, subreddits, args.sleep, args.limit


async def main():
    mode, subreddits, sleep_interval, limit = parse_args()

    scrapers = [
        RedditScraper(subreddit, mode, limit=limit) for subreddit in subreddits
    ]

    while True:
        for subreddit, scraper in zip(subreddits, scrapers):
            logger.info(f'Scraping for new {mode} on r/{subreddit} at {datetime.now()}...')

            try:
                num_results = await scraper.run()
            except RequestException:
                logger.info(f'Error scraping on r/{subreddit}')
                continue

            logger.info(f'Saved {num_results} {mode}')

        gc.collect()
        logger.info(f'Sleeping for {sleep_interval} seconds...')
        sleep(sleep_interval)


if __name__ == '__main__':
    asyncio.run(main())