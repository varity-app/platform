"""
Unit tests for the reddit submissions pipeline
"""

from typing import List, Dict
import pytest

from util.constants.scraping import (
    DataSources as DS,
    ParentSources as PS,
    ScrapedPostConstants as SPC,
)

from ..submissions import parse_scraped_posts

test_data = [
    (
        {
            "submission_id": "mpudu8",
            "subreddit": "wallstreetbets",
            "title": "My Options account just got converted into diamond hand ape bag"
            " 100% IN...Almost. original purchase in Dec. I am finally at peace.",
            "created_utc": "2021-04-13T03:42:33",
            "name": "t3_mpudu8",
            "selftext": "No selftext to go along here...",
            "author": "KoreanSpyPuts",
            "is_original_content": False,
            "is_text": False,
            "nsfw": False,
            "num_comments": 20,
            "permalink": "/r/wallstreetbets/comments/mpudu8/my_options_"
            "account_just_got_converted_into/",
            "upvotes": 269,
            "upvote_ratio": 0.91,
            "url": "https://i.redd.it/qay238d93vs61.png",
        },
        [
            {
                SPC.PARENT_ID: "mpudu8",
                SPC.PARENT_SOURCE: PS.SUBMISSION_TITLE,
                SPC.DATA_SOURCE: DS.REDDIT,
                SPC.TEXT: "My Options account just got converted into diamond hand ape"
                " bag 100% IN...Almost. original purchase in Dec. I am finally at peace.",
                SPC.TIMESTAMP: "2021-04-13T03:42:33",
            },
            {
                SPC.PARENT_ID: "mpudu8",
                SPC.PARENT_SOURCE: PS.SUBMISSION_SELFTEXT,
                SPC.DATA_SOURCE: DS.REDDIT,
                SPC.TEXT: "No selftext to go along here...",
                SPC.TIMESTAMP: "2021-04-13T03:42:33",
            },
        ],
    ),
    (
        {
            "submission_id": "mpudu8",
            "subreddit": "wallstreetbets",
            "title": "My Options account just got converted into diamond hand ape"
            " bag 100% IN...Almost. original purchase in Dec. I am finally at peace.",
            "created_utc": "2021-04-13T03:42:33",
            "name": "t3_mpudu8",
            "selftext": "",
            "author": "KoreanSpyPuts",
            "is_original_content": False,
            "is_text": False,
            "nsfw": False,
            "num_comments": 20,
            "permalink": "/r/wallstreetbets/comments/mpudu8/my_options"
            "_account_just_got_converted_into/",
            "upvotes": 269,
            "upvote_ratio": 0.91,
            "url": "https://i.redd.it/qay238d93vs61.png",
        },
        [
            {
                SPC.PARENT_ID: "mpudu8",
                SPC.PARENT_SOURCE: PS.SUBMISSION_TITLE,
                SPC.DATA_SOURCE: DS.REDDIT,
                SPC.TEXT: "My Options account just got converted into diamond hand ape"
                " bag 100% IN...Almost. original purchase in Dec. I am finally at peace.",
                SPC.TIMESTAMP: "2021-04-13T03:42:33",
            }
        ],
    ),
]


@pytest.mark.parametrize("submission,scraped_posts", test_data)
def test_pipeline(submission: Dict, scraped_posts: List[Dict]) -> None:
    """Unit test the submissions parse_scraped_posts extraction method"""
    parsed_posts = parse_scraped_posts(submission)

    match_count = 0
    for post in parsed_posts:
        for other_post in scraped_posts:
            if post == other_post:
                match_count += 1

    assert len(parsed_posts) == len(scraped_posts)
    assert match_count == len(scraped_posts)
