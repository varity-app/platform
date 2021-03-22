import faust

from util.tickers import parse_tickers, all_tickers
from util.constants.scraping import DataSources as DS, ParentSources as PS
from util.constants import MentionTypes

from process.app import app
from .views import comments_topic
from .models import Comment
from ..ticker_mentions.views import ticker_mentions_topic
from ..ticker_mentions.models import TickerMention
from ..scraped_posts.views import scraped_posts_topic
from ..scraped_posts.models import ScrapedPost


async def parse_ticker_fields(comment: Comment) -> bool:
    """Parse tickers from body to TickerMention messages and publish to Kafka"""
    # Parse Tickers
    body_tickers = parse_tickers(
        comment.body, all_tickers=all_tickers
    )

    # Publish to Kafka
    for ticker in body_tickers:
        mention = TickerMention(
            ticker,
            DS.REDDIT,
            PS.COMMENT_BODY,
            comment.comment_id,
            comment.created_utc,
            MentionTypes.TICKER,
        )
        await ticker_mentions_topic.send(value=mention)

    # Return whether post has any tickers mentioned
    has_tickers = len(body_tickers) > 0
    return has_tickers


async def parse_post(comment: Comment):
    """Parse body field to ScrapedPost message and publish to Kafka"""
    # Publish body
    body_post = ScrapedPost(
        comment.body,
        DS.REDDIT,
        PS.COMMENT_BODY,
        comment.comment_id,
        comment.created_utc,
    )
    await scraped_posts_topic.send(value=body_post)


@app.agent(comments_topic)
async def process_comment(comments):
    async for comment in comments:
        has_tickers = await parse_ticker_fields(comment)
        if has_tickers:
            await parse_post(comment)
