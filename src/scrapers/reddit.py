"""
Define the RedditScraper class
"""

from datetime import datetime
from typing import Union
import logging

from confluent_kafka import Producer
import asyncpraw

from util.constants.reddit import (
    CommentConstants as CC,
    SubmissionConstants as SC,
    Config as RC,
    Misc,
)
from util.constants.kafka import Config, Topics
from util.constants.dynamo import Tables as DT
from util.dynamo import DynamoTable
from messages.reddit.submission import SubmissionMessage
from messages.reddit.comment import CommentMessage

from . import Scraper

logger = logging.getLogger(__name__)


class RedditScraper(Scraper):
    """
    The scraper checks for new posts at a specified interval, and if they
    have not been previously processed, submits them to Kafka.

    The scraper's memory consists of a local Memory helper class and a DynamoDB table.
    """

    def __init__(
        self,
        subreddit: str,
        mode: str,
        limit=100,
        enable_publish=True,
        enable_dynamo=True,
    ) -> None:
        super().__init__(limit)

        self.subreddit = subreddit
        self.limit = limit
        self.enable_publish = enable_publish
        self.enable_dynamo = enable_dynamo

        assert mode in [Misc.COMMENTS, Misc.SUBMISSIONS]
        self.mode = mode

        self.reddit = asyncpraw.Reddit(
            client_id=RC.REDDIT_CLIENT_ID,
            client_secret=RC.REDDIT_CLIENT_SECRET,
            user_agent=RC.REDDIT_USER_AGENT,
            username=RC.REDDIT_USERNAME,
            password=RC.REDDIT_PASSWORD,
        )

        if enable_publish:
            self.publisher = Producer(Config.OBJ)

        if enable_dynamo:
            if mode == Misc.SUBMISSIONS:
                self.dynamo = DynamoTable(DT.REDDIT_SUBMISSIONS, SC.ID)
            else:
                self.dynamo = DynamoTable(DT.REDDIT_COMMENTS, CC.ID)

    @staticmethod
    def callback(err, msg):
        """Kafka publisher callback"""
        if err is not None:
            logger.error(f"Failed to deliver message: {err}")
        else:
            logger.debug(
                f"Produced record to topic {msg.topic()}"
                "partition [{msg.partition()}] @ offset {msg.offset()}"
            )

    def publish(
        self, topic: str, message: Union[SubmissionMessage, CommentMessage]
    ) -> int:
        """Serialize a message and publish it to Kafka"""
        self.publisher.produce(topic, value=message.serialize())

    async def scrape_comments(self):
        """Scrape newest comments from Reddit"""

        subreddit_origin = await self.reddit.subreddit(self.subreddit)

        comment_count = 0
        async for comment in subreddit_origin.comments(limit=self.limit):
            if self.memory_contains(comment.id):
                continue

            if self.enable_dynamo and self.dynamo.exists(comment.id):
                continue

            self.memory_add(comment.id)

            # Parse Comment
            comment = self.parse_comment(comment)

            # Save in Kafka
            if self.enable_publish:
                self.publish(Topics.REDDIT_COMMENTS, comment)

            # Save in DynamoDB
            if self.enable_dynamo:
                self.dynamo.put(comment.comment_id)

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
            if self.memory_contains(submission.id):
                continue

            if self.enable_dynamo and self.dynamo.exists(submission.id):
                continue

            self.memory_add(submission.id)

            # Parse Submission
            submission = self.parse_submission(submission)

            # Save in Kafka
            if self.enable_publish:
                self.publish(Topics.REDDIT_SUBMISSIONS, submission)

            # Save in DynamoDB
            if self.enable_dynamo:
                self.dynamo.put(submission.submission_id)

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
        self.memory_reset_cursor()

        if self.mode == "submissions":
            num_results = await self.scrape_submissions()
        else:
            num_results = await self.scrape_comments()

        self.publisher.flush()

        return num_results
