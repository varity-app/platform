"""
Define the RedditScraper class
"""

from datetime import datetime
from typing import Union
import logging

from google.cloud import pubsub_v1
import asyncpraw

from util.constants.reddit import (
    Config as RC,
    Misc,
)
from util.constants.pubsub import Config, Topics
from util.constants.firestore import Collections
from messages.reddit.submission import SubmissionMessage
from messages.reddit.comment import CommentMessage

from .memory import Memory

logger = logging.getLogger(__name__)


class RedditScraper:
    """
    The scraper checks for new posts at a specified interval, and if they
    have not been previously processed, submits them to Pub/Sub.

    The scraper's memory consists of a local Memory helper class and a DynamoDB table.
    """

    def __init__(
        self,
        subreddit: str,
        mode: str,
        limit=100,
        enable_publish=True,
        enable_firestore=True,
    ) -> None:
        assert mode in [Misc.COMMENTS, Misc.SUBMISSIONS]

        collection = Collections.REDDIT_COMMENTS
        if mode == Misc.SUBMISSIONS:
            collection = Collections.REDDIT_SUBMISSIONS

        self.memory = Memory(limit, enable_firestore, collection=collection)

        self.subreddit = subreddit
        self.limit = limit
        self.enable_publish = enable_publish

        self.mode = mode

        self.reddit = asyncpraw.Reddit(
            client_id=RC.REDDIT_CLIENT_ID,
            client_secret=RC.REDDIT_CLIENT_SECRET,
            user_agent=RC.REDDIT_USER_AGENT,
            username=RC.REDDIT_USERNAME,
            password=RC.REDDIT_PASSWORD,
        )

        # Initialize Pub/Sub publisher
        if enable_publish:
            self.publisher = pubsub_v1.PublisherClient()
            if mode == Misc.SUBMISSIONS:
                self.topic = self.publisher.topic_path(  # pylint: disable=no-member
                    Config.PROJECT, Topics.REDDIT_SUBMISSIONS
                )
            else:
                self.topic = self.publisher.topic_path(  # pylint: disable=no-member
                    Config.PROJECT, Topics.REDDIT_COMMENTS
                )

    def publish(self, message: Union[SubmissionMessage, CommentMessage]) -> int:
        """Serialize a message and publish it to Pub/Sub"""
        self.publisher.publish(self.topic, message.serialize().encode("utf-8")).result()

    async def scrape_comments(self):
        """Scrape newest comments from Reddit"""

        subreddit_origin = await self.reddit.subreddit(self.subreddit)

        comment_count = 0
        async for comment in subreddit_origin.comments(limit=self.limit):
            if self.memory.contains(comment.id):
                continue

            self.memory.add(comment.id)

            # Parse Comment
            comment = self.parse_comment(comment)

            # Save in Pub/Sub
            if self.enable_publish:
                self.publish(comment)

            comment_count += 1

        return comment_count

    def parse_comment(self, comment):
        """Parse a PRAW comment to the custom CommentMessage class"""
        created_utc = datetime.utcfromtimestamp(comment.created_utc).isoformat()

        if comment.author is not None:
            author = comment.author.name
        else:
            author = ""

        com_obj = CommentMessage(
            comment.id,
            comment.link_id,
            self.subreddit,
            author,
            created_utc,
            comment.body,
            comment.score,
        )

        return com_obj

    async def scrape_submissions(self):
        """Scrape newest submissions from Reddit"""
        subreddit_origin = await self.reddit.subreddit(self.subreddit)

        submission_count = 0
        async for submission in subreddit_origin.new(limit=self.limit):
            if self.memory.contains(submission.id):
                continue

            self.memory.add(submission.id)

            # Parse Submission
            submission = self.parse_submission(submission)

            # Save in Pub/Sub
            if self.enable_publish:
                self.publish(submission)

            submission_count += 1

        return submission_count

    def parse_submission(self, submission):
        """Parse a PRAW submission to the custom SubmissionMessage class"""
        submission_id = submission.id
        title = submission.title
        created_utc = datetime.utcfromtimestamp(submission.created_utc).isoformat()
        is_original_content = submission.is_original_content
        is_text = submission.is_self
        name = submission.name
        num_comments = submission.num_comments
        nsfw = submission.over_18
        permalink = submission.permalink
        upvotes = submission.score
        selftext = submission.selftext
        upvote_ratio = submission.upvote_ratio
        url = submission.url
        if submission.author is not None:
            author = submission.author.name
        else:
            author = ""

        sub_obj = SubmissionMessage(
            submission_id,
            self.subreddit,
            title,
            created_utc,
            name,
            selftext,
            author,
            is_original_content,
            is_text,
            nsfw,
            num_comments,
            permalink,
            upvotes,
            upvote_ratio,
            url,
        )

        return sub_obj

    async def run(self):
        """Entrypoint method"""
        self.memory.reset_cursor()

        if self.mode == "submissions":
            num_results = await self.scrape_submissions()
        else:
            num_results = await self.scrape_comments()

        return num_results
