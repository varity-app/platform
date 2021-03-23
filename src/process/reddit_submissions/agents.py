"""
Faust Agents and helper methods for processing Reddit Submissions
"""

from typing import List, Tuple

from util.tickers import parse_tickers, all_tickers
from util.constants.scraping import DataSources as DS, ParentSources as PS
from util.constants import MentionTypes

from process.app import app
from .views import submissions_topic
from .models import Submission
from ..ticker_mentions.views import ticker_mentions_topic
from ..ticker_mentions.models import TickerMention
from ..scraped_posts.views import scraped_posts_topic
from ..scraped_posts.models import ScrapedPost


def parse_ticker_fields(submission: Submission) -> Tuple[List[str], List[str]]:
    """Parse tickers from selftext and title fields to TickerMention messages"""
    # Parse Tickers
    selftext_tickers = parse_tickers(submission.selftext, all_tickers=all_tickers)

    title_tickers = parse_tickers(submission.title, all_tickers=all_tickers)

    return selftext_tickers, title_tickers


def create_mention_object(
    ticker: str, submission: Submission, parent_source: str
) -> TickerMention:
    """Create a TickerMention object for a ticker"""
    mention = TickerMention(
        stock_name=ticker,
        data_source=DS.REDDIT,
        parent_source=parent_source,
        parent_id=submission.submission_id,
        created_utc=submission.created_utc,
        mention_type=MentionTypes.TICKER,
    )

    return mention


def parse_posts(submission: Submission) -> Tuple[ScrapedPost, ScrapedPost]:
    """Parse selftext and title fields to ScrapedPost messages"""
    # Publish selftext
    selftext_post = ScrapedPost(
        text=submission.selftext,
        data_source=DS.REDDIT,
        parent_source=PS.SUBMISSION_SELFTEXT,
        parent_id=submission.submission_id,
        timestamp=submission.created_utc,
    )

    # Parse title
    title_post = ScrapedPost(
        text=submission.title,
        data_source=DS.REDDIT,
        parent_source=PS.SUBMISSION_TITLE,
        parent_id=submission.submission_id,
        timestamp=submission.created_utc,
    )

    return selftext_post, title_post


@app.agent(submissions_topic)
async def process_submission(submissions):
    """Parse tickers and ScrapedPosts from Reddit Submissions"""

    async for submission in submissions:
        # Parse tickers
        selftext_tickers, title_tickers = parse_ticker_fields(submission)

        # Publish to Kafka
        for ticker in selftext_tickers:
            mention = create_mention_object(ticker, submission, PS.SUBMISSION_SELFTEXT)
            await ticker_mentions_topic.send(value=mention)

        for ticker in title_tickers:
            mention = create_mention_object(ticker, submission, PS.SUBMISSION_TITLE)
            await ticker_mentions_topic.send(value=mention)

        # Publish ScrapedPost if necessary
        has_tickers = len(selftext_tickers + title_tickers) > 0

        if has_tickers:
            selftext_post, title_post = parse_posts(submission)

            if submission.selftext:
                await scraped_posts_topic.send(value=selftext_post)

            await scraped_posts_topic.send(value=title_post)
