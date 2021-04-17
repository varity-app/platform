"""
Unit tests for the reddit comments pipeline
"""

from typing import Dict
import pytest

from util.constants.scraping import (
    DataSources as DS,
    ParentSources as PS,
    ScrapedPostConstants as SPC,
)

from ..comments import parse_scraped_post

test_data = [
    (
        {
            "comment_id": "gug9773",
            "submission_id": "t3_mqaa3x",
            "subreddit": "wallstreetbets",
            "author": "ricardosanch5",
            "created_utc": "2021-04-14T03:55:18",
            "body": "10, and the person lost their time",
            "upvotes": 1,
        },
        {
            SPC.PARENT_ID: "gug9773",
            SPC.PARENT_SOURCE: PS.COMMENT_BODY,
            SPC.DATA_SOURCE: DS.REDDIT,
            SPC.TEXT: "10, and the person lost their time",
            SPC.TIMESTAMP: "2021-04-14T03:55:18",
        },
    ),
]


@pytest.mark.parametrize("comment,scraped_post", test_data)
def test_pipeline(comment: Dict, scraped_post: Dict) -> None:
    """Unit test the comments parse_scraped_post extraction method"""
    parsed_post = parse_scraped_post(comment)

    assert parsed_post == scraped_post
