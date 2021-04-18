"""
Define the scraper class for scraping historical reddit posts with the psaw library
"""

from datetime import datetime
from typing import Union
import logging

from google.cloud import pubsub_v1
import psaw

from util.constants.reddit import (
    Misc,
)
from util.constants.pubsub import Config, Topics
from util.constants.firestore import Collections
from messages.reddit.submission import SubmissionMessage
from messages.reddit.comment import CommentMessage

from ..memory import Memory

logger = logging.getLogger(__name__)


class HistoricalRedditScraper:
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
        filter_removed=True,
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
        self.filter_removed = filter_removed

        self.mode = mode

        self.api = psaw.PushshiftAPI()

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

    @staticmethod
    def was_removed(submission) -> bool:
        """Check if a submission was removed"""
        try:
            was_removed = submission.selftext in ["[removed]", "[delted]"]
        except AttributeError:
            return True

        return was_removed

    def scrape_comments(self, start_date: datetime, end_date: datetime):
        """Scrape historical comments from Reddit"""
        start_date, end_date = int(start_date.timestamp()), int(end_date.timestamp())

        gen = self.api.search_comments(
            subreddit=self.subreddit,
            limit=self.limit,
            after=start_date,
            before=end_date,
        )

        comment_count = 0
        for comment in gen:
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
        """Parse a PSAW comment to the custom CommentMessage class"""
        created_utc = datetime.utcfromtimestamp(comment.created_utc).isoformat()

        if comment.author is not None:
            author = comment.author
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

    def scrape_submissions(self, start_date: datetime, end_date: datetime) -> None:
        """Scrape past submissions from Reddit"""
        start_date, end_date = int(start_date.timestamp()), int(end_date.timestamp())
        gen = self.api.search_submissions(
            subreddit=self.subreddit,
            limit=self.limit,
            after=start_date,
            before=end_date,
        )

        submission_count = 0
        for submission in gen:
            if self.filter_removed and self.was_removed(submission):
                continue

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
        """Parse a PSAW submission to the custom SubmissionMessage class"""
        submission_id = submission.id
        title = submission.title
        created_utc = datetime.utcfromtimestamp(submission.created_utc).isoformat()

        try:
            is_original_content = submission.is_original_content
        except AttributeError:
            is_original_content = False

        is_text = submission.is_self
        name = submission.id  # PSAW has no name field
        num_comments = submission.num_comments
        nsfw = submission.over_18
        permalink = submission.permalink
        upvotes = submission.score

        try:
            selftext = submission.selftext
        except AttributeError:
            selftext = ""

        upvote_ratio = 0  # PSAW has no upvote_ratio field
        url = submission.url
        if submission.author is not None:
            author = submission.author
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

    def run(self, start_date: datetime, end_date: datetime) -> None:
        """Entrypoint method"""

        if self.mode == "submissions":
            num_results = self.scrape_submissions(start_date, end_date)
        else:
            num_results = self.scrape_comments(start_date, end_date)

        return num_results
