import asyncpraw
from datetime import datetime
import boto3

from util.constants.reddit import Config as RC, CommentConstants as CC, SubmissionConstants as SC, Misc
from util.constants.aws import SNS
from messages.reddit.submission import SubmissionMessage
from messages.reddit.comment import CommentMessage
from messages.sns import SNSMessage
from . import Scraper


class RedditScraper(Scraper):
    def __init__(
        self, subreddit, mode, limit=100, enable_publish=True,
    ):
        super().__init__(limit)

        self.subreddit = subreddit
        self.limit = limit
        self.enable_publish = enable_publish

        assert mode in [Misc.COMMENTS, Misc.SUBMISSIONS]
        self.mode = mode

        self.reddit = asyncpraw.Reddit(
            client_id=RC.REDDIT_CLIENT_ID,
            client_secret=RC.REDDIT_CLIENT_SECRET,
            user_agent=RC.REDDIT_USER_AGENT,
            username=RC.REDDIT_USERNAME,
            password=RC.REDDIT_PASSWORD,
        )

        self.sns = boto3.client('sns')

    def publish(self, arn, message):
        """Serialize a message and publish it to AWS SNS"""
        # response = self.sns.publish(
        #     TargetArn=arn,
        #     Message=SNSMessage.serialize(message),
        #     MessageStructure=SNS.JSON,
        # )
        pass

    async def scrape_comments(self):
        """Scrape newest comments from Reddit"""

        subreddit_origin = await self.reddit.subreddit(self.subreddit)

        comment_count = 0
        async for comment in subreddit_origin.comments(limit=self.limit):
            if self.memory_contains(comment.id):
                continue

            self.memory_add(comment.id)

            # Save Comment
            comment = self.parse_comment(comment)
            if self.enable_publish:
                self.publish(SNS.REDDIT_COMMENTS, comment)

            comment_count += 1

        return comment_count

    def parse_comment(self, comment):
        created_utc = datetime.utcfromtimestamp(
            comment.created_utc).isoformat()

        if comment.author is not None:
            author = comment.author.name
        else:
            author = ''

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

            self.memory_add(submission.id)

            # Save Submission
            submission = self.parse_submission(submission)
            if self.enable_publish:
                self.publish(SNS.REDDIT_SUBMISSIONS, submission)

            submission_count += 1

        return submission_count

    def parse_submission(self, submission):
        submission_id = submission.id
        title = submission.title
        created_utc = datetime.utcfromtimestamp(
            submission.created_utc).isoformat()
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
            author = ''

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

        return num_results
