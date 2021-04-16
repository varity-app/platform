"""
Unit tests for the scraped posts pipeline
"""

from typing import List, Dict
import pytest


import apache_beam as beam

from util.constants.scraping import (
    TickerMentionsConstants as TMC,
    MentionTypes,
)

from ..scraped_posts import parse_ticker_mentions

test_data = [
    (
        {
            "parent_id": "mpuizf",
            "timestamp": "2021-04-13T03:51:58",
            "text": "Moderate Gain from $BABA",
            "data_source": "reddit",
            "parent_source": "submission_title",
        },
        [
            {
                TMC.TICKER: "BABA",
                TMC.DATA_SOURCE: "reddit",
                TMC.PARENT_SOURCE: "submission_title",
                TMC.PARENT_ID: "mpuizf",
                TMC.TIMESTAMP: "2021-04-13T03:51:58",
                TMC.MENTION_TYPE: MentionTypes.TICKER,
            }
        ],
    ),
    (
        {
            "parent_id": "mpuizf",
            "timestamp": "2021-04-13T03:51:58",
            "text": "Was averaging down and holding/rolling since Jan. Hopefully can see $BABA @ 300+ EOY.",
            "data_source": "reddit",
            "parent_source": "submission_selftext",
        },
        [
            {
                TMC.TICKER: "BABA",
                TMC.DATA_SOURCE: "reddit",
                TMC.PARENT_SOURCE: "submission_selftext",
                TMC.PARENT_ID: "mpuizf",
                TMC.TIMESTAMP: "2021-04-13T03:51:58",
                TMC.MENTION_TYPE: MentionTypes.TICKER,
            }
        ],
    ),
    (
        {
            "parent_id": "mpudu8",
            "timestamp": "2021-04-13T03:42:33",
            "text": "My Options account just got converted into diamond hand ape bag 100% IN...Almost. original purchase in Dec. I am finally at peace.",
            "data_source": "reddit",
            "parent_source": "submission_title",
        },
        [],
    ),
]


@pytest.mark.parametrize("scraped_post,ticker_mentions", test_data)
def test_pipeline(scraped_post: Dict, ticker_mentions: List[Dict]) -> None:
    """Unit test the scraped posts parse_ticker_mentions extraction method"""
    parsed_mentions = parse_ticker_mentions(scraped_post)

    match_count = 0
    for mentions in parsed_mentions:
        for other_mentions in ticker_mentions:
            if mentions == other_mentions:
                match_count += 1

    assert len(parsed_mentions) == len(ticker_mentions)
    assert match_count == len(ticker_mentions)
