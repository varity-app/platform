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


def parse_ticker_fields(comment: Comment) -> bool:
    """Parse tickers from body to TickerMention messages and publish to Kafka"""
    body_tickers = parse_tickers(
        comment.body, all_tickers=all_tickers
    )

    return body_tickers


def create_mention_object(ticker: str, comment: Comment) -> TickerMention:
    mention = TickerMention(
        stock_name=ticker,
        data_source=DS.REDDIT,
        parent_source=PS.COMMENT_BODY,
        parent_id=comment.comment_id,
        created_utc=comment.created_utc,
        mention_type=MentionTypes.TICKER,
    )

    return mention


def parse_post(comment: Comment) -> None:
    """Parse body field to ScrapedPost message"""
    body_post = ScrapedPost(
        text=comment.body,
        data_source=DS.REDDIT,
        parent_source=PS.COMMENT_BODY,
        parent_id=comment.comment_id,
        timestamp=comment.created_utc,
    )
    
    return body_post


@app.agent(comments_topic)
async def process_comment(comments) -> None:
    """Parse a reddit comment for tickers and publish as a ScrapedPost"""

    async for comment in comments:
        # Parse Tickers
        body_tickers = parse_ticker_fields(comment)
        for ticker in body_tickers:
            ticker_mention = create_mention_object(ticker, comment)
            await ticker_mentions_topic.send(value=mention)

        # Publish as ScrapedPost if there are tickers
        has_tickers = len(body_tickers) > 0

        if has_tickers:
            body_post = parse_post(comment)
            await scraped_posts_topic.send(value=body_post)
