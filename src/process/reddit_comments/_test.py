"""
Unit tests for processing Reddit Comments from Kafka
"""

import pytest

from util.constants.scraping import DataSources as DS, ParentSources as PS
from util.constants import MentionTypes

from .agents import parse_ticker_fields, create_mention_object, parse_post
from .models import Comment


@pytest.fixture
def comment_with_no_ticker():
    """A comment with no tickers"""
    inp = Comment(
        comment_id="grrmy6c",
        submission_id="t3_maauj7",
        subreddit="wallstreetbets",
        author="DaySwingTrade",
        created_utc="2021-03-23T03:21:46",
        body="No tickers here!  Just a big Long Sentence...",
        upvotes=48,
    )
    return inp


@pytest.fixture
def comment_with_one_ticker():
    """A comment with one ticker"""
    inp = Comment(
        comment_id="grrmy6c",
        submission_id="t3_maauj7",
        subreddit="wallstreetbets",
        author="DaySwingTrade",
        created_utc="2021-03-23T03:21:46",
        body="I just realized more people going to people tuning into GME earnings",
        upvotes=48,
    )
    return inp


@pytest.fixture
def comment_with_many_tickers():
    """A comment with many tickers"""
    inp = Comment(
        comment_id="grrmy6c",
        submission_id="t3_maauj7",
        subreddit="wallstreetbets",
        author="DaySwingTrade",
        created_utc="2021-03-23T03:21:46",
        body="Stonk spam bot 123: $GME $GME $NOK $BB $TSLA AAPL... Buy all!",
        upvotes=48,
    )
    return inp


def test_parse_tickers_no_ticker(comment_with_no_ticker):
    """Test ticker parsing for comment with no tickers"""
    comment = comment_with_no_ticker
    body_tickers = parse_ticker_fields(comment)

    assert body_tickers == []


def test_parse_tickers_one_ticker(comment_with_one_ticker):
    """Test ticker parsing for comment with one ticker"""
    comment = comment_with_one_ticker
    body_tickers = parse_ticker_fields(comment)

    assert body_tickers == ["GME"]

    for ticker in body_tickers:
        ticker_mention = create_mention_object(ticker, comment)

        assert ticker_mention.stock_name == ticker
        assert ticker_mention.parent_id == comment.comment_id
        assert ticker_mention.data_source == DS.REDDIT
        assert ticker_mention.parent_source == PS.COMMENT_BODY
        assert ticker_mention.created_utc == comment.created_utc
        assert ticker_mention.mention_type == MentionTypes.TICKER


def test_parse_tickers_many_tickers(comment_with_many_tickers):
    """Test ticker parsing for comment with many tickers"""
    comment = comment_with_many_tickers
    body_tickers = parse_ticker_fields(comment)

    correct_tickers = ["GME", "NOK", "BB", "TSLA", "AAPL"]

    print(body_tickers, correct_tickers)

    assert sorted(body_tickers) == sorted(correct_tickers)

    for ticker in body_tickers:
        ticker_mention = create_mention_object(ticker, comment)

        assert ticker_mention.stock_name == ticker
        assert ticker_mention.parent_id == comment.comment_id
        assert ticker_mention.data_source == DS.REDDIT
        assert ticker_mention.parent_source == PS.COMMENT_BODY
        assert ticker_mention.created_utc == comment.created_utc
        assert ticker_mention.mention_type == MentionTypes.TICKER


def test_parse_post_no_ticker(comment_with_no_ticker):
    """Test ScrapedPost parsing for comment with no tickers"""
    comment = comment_with_no_ticker
    post = parse_post(comment)

    assert post.parent_id == comment.comment_id
    assert post.data_source == DS.REDDIT
    assert post.parent_source == PS.COMMENT_BODY
    assert post.text == comment.body
    assert post.timestamp == comment.created_utc


def test_parse_post_one_ticker(comment_with_one_ticker):
    """Test ScrapedPost parsing for comment with one ticker"""
    comment = comment_with_one_ticker
    post = parse_post(comment)

    assert post.parent_id == comment.comment_id
    assert post.data_source == DS.REDDIT
    assert post.parent_source == PS.COMMENT_BODY
    assert post.text == comment.body
    assert post.timestamp == comment.created_utc
