import faust

from util.tickers import parse_tickers, all_tickers
from util.constants.scraping import DataSources as DS, ParentSources as PS
from util.constants import MentionTypes

from process.app import app
from .views import submissions_topic
from ..ticker_mentions.views import ticker_mentions_topic
from ..ticker_mentions.models import TickerMention
from ..scraped_posts.views import scraped_posts_topic
from ..scraped_posts.models import ScrapedPost


async def parse_ticker_fields(submission) -> bool:
    """Parse tickers from selftext and title fields to TickerMention messages and publish them to Kafka"""
    # Parse Tickers
    selftext_tickers = parse_tickers(
        submission.selftext, all_tickers=all_tickers
    )

    title_tickers = parse_tickers(
        submission.title, all_tickers=all_tickers)

    # Publish to Kafka
    for ticker in selftext_tickers:
        mention = TickerMention(
            stock_name=ticker,
            data_source=DS.REDDIT,
            parent_source=PS.SUBMISSION_SELFTEXT,
            parent_id=submission.submission_id,
            created_utc=submission.created_utc,
            mention_type=MentionTypes.TICKER,
        )
        await ticker_mentions_topic.send(value=mention)

    for ticker in title_tickers:
        mention = TickerMention(
            stock_name=ticker,
            data_source=DS.REDDIT,
            parent_source=PS.SUBMISSION_TITLE,
            parent_id=submission.submission_id,
            created_utc=submission.created_utc,
            mention_type=MentionTypes.TICKER,
        )
        await ticker_mentions_topic.send(value=mention)

    # Return whether post has any tickers mentioned
    has_tickers = len(selftext_tickers) > 0 or len(title_tickers) > 0
    return has_tickers


async def parse_posts(submission):
    """Parse selftext and title fields to ScrapedPost messages and publish them to Kafka"""
    # Publish selftext if necessary
    if submission.selftext:
        selftext_post = ScrapedPost(
            text=submission.selftext,
            data_source=DS.REDDIT,
            parent_source=PS.SUBMISSION_SELFTEXT,
            parent_id=submission.submission_id,
            timestamp=submission.created_utc,
        )
        await scraped_posts_topic.send(value=selftext_post)

    # Publish title
    title_post = ScrapedPost(
        text=submission.title,
        data_source=DS.REDDIT,
        parent_source=PS.SUBMISSION_TITLE,
        parent_id=submission.submission_id,
        timestamp=submission.created_utc,
    )
    await scraped_posts_topic.send(value=title_post)


@app.agent(submissions_topic)
async def process_submission(submissions):
    async for submission in submissions:
        has_tickers = await parse_ticker_fields(submission)
        if has_tickers:
            await parse_posts(submission)
