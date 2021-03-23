import pytest

from util.constants.scraping import DataSources as DS, ParentSources as PS
from util.constants import MentionTypes

from .agents import parse_ticker_fields, create_mention_object, parse_posts
from .models import Submission
from ..scraped_posts.models import ScrapedPost
from ..ticker_mentions.models import TickerMention


@pytest.fixture
def submission_with_no_ticker():
    inp = Submission(
        submission_id="m8g2ec",
        subreddit="smallstreetbets",
        title="Why do I even post?",
        created_utc="2021-03-19T12:23:44",
        name="t3_m8g2ec",
        selftext="Nothing interesting here...",
        author="superlopster",
        is_original_content=False,
        is_text=True,
        nsfw=False,
        num_comments=14,
        permalink="/r/smallstreetbets/comments/m8g2ec/every_thing_is_fine/",
        upvotes=2,
        upvote_ratio=0.58,
        url="https://www.reddit.com/r/smallstreetbets/comments/m8g2ec/every_thing_is_fine/"
    )
    return inp


@pytest.fixture
def submission_with_one_ticker():
    inp = Submission(
        submission_id="m8g2ec",
        subreddit="smallstreetbets",
        title="Every thing is fine .",
        created_utc="2021-03-19T12:23:44",
        name="t3_m8g2ec",
        selftext="Today RBLX fall 12.36 % and I started to think why . Opent up news and evry thing was positive good news, all is fine but stock falls 😂",
        author="superlopster",
        is_original_content=False,
        is_text=True,
        nsfw=False,
        num_comments=14,
        permalink="/r/smallstreetbets/comments/m8g2ec/every_thing_is_fine/",
        upvotes=2,
        upvote_ratio=0.58,
        url="https://www.reddit.com/r/smallstreetbets/comments/m8g2ec/every_thing_is_fine/"
    )
    return inp


@pytest.fixture
def submission_with_multiple_tickers():
    inp = Submission(
        submission_id="m8g2ec",
        subreddit="smallstreetbets",
        title="GME to the Moon",
        created_utc="2021-03-19T12:23:44",
        name="t3_m8g2ec",
        selftext="PLTR is doing better then TSLA imo.  Still, blue chips like AAPL are king",
        author="superlopster",
        is_original_content=False,
        is_text=True,
        nsfw=False,
        num_comments=14,
        permalink="/r/smallstreetbets/comments/m8g2ec/every_thing_is_fine/",
        upvotes=2,
        upvote_ratio=0.58,
        url="https://www.reddit.com/r/smallstreetbets/comments/m8g2ec/every_thing_is_fine/"
    )
    return inp


def test_parse_tickers_no_ticker(submission_with_no_ticker):
    submission = submission_with_no_ticker
    selftext_tickers, title_tickers = parse_ticker_fields(submission)

    assert selftext_tickers == []
    assert title_tickers == []


def test_parse_tickers_one_ticker(submission_with_one_ticker):
    submission = submission_with_one_ticker
    selftext_tickers, title_tickers = parse_ticker_fields(submission)

    assert selftext_tickers == ["RBLX"]
    assert title_tickers == []

    for ticker in selftext_tickers:
        ticker_mention = create_mention_object(ticker, submission, PS.SUBMISSION_SELFTEXT)

        assert ticker_mention.stock_name == ticker
        assert ticker_mention.parent_id == submission.submission_id
        assert ticker_mention.data_source == DS.REDDIT
        assert ticker_mention.parent_source == PS.SUBMISSION_SELFTEXT
        assert ticker_mention.created_utc == submission.created_utc
        assert ticker_mention.mention_type == MentionTypes.TICKER


def test_parse_tickers_multiple_tickers(submission_with_multiple_tickers):
    submission = submission_with_multiple_tickers
    selftext_tickers, title_tickers = parse_ticker_fields(submission)

    assert sorted(selftext_tickers) == sorted(["PLTR", "TSLA", "AAPL"])
    assert title_tickers == ["GME"]

    for ticker in selftext_tickers:
        ticker_mention = create_mention_object(ticker, submission, PS.SUBMISSION_SELFTEXT)

        assert ticker_mention.stock_name == ticker
        assert ticker_mention.parent_id == submission.submission_id
        assert ticker_mention.data_source == DS.REDDIT
        assert ticker_mention.parent_source == PS.SUBMISSION_SELFTEXT
        assert ticker_mention.created_utc == submission.created_utc
        assert ticker_mention.mention_type == MentionTypes.TICKER

    for ticker in title_tickers:
        ticker_mention = create_mention_object(ticker, submission, PS.SUBMISSION_TITLE)

        assert ticker_mention.stock_name == ticker
        assert ticker_mention.parent_id == submission.submission_id
        assert ticker_mention.data_source == DS.REDDIT
        assert ticker_mention.parent_source == PS.SUBMISSION_TITLE
        assert ticker_mention.created_utc == submission.created_utc
        assert ticker_mention.mention_type == MentionTypes.TICKER