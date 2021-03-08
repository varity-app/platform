import boto3

from messages.reddit.submission import SubmissionMessage
from messages.scraped_post import ScrapedPostMessage
from messages.ticker_mention import TickerMentionMessage
from messages.sns import SNSMessage

from util.tickers import parse_tickers, all_tickers
from util.constants.scraping import DataSources as DS, ParentSources as PS
from util.constants.aws import SNS


class SubmissionConsumer(object):
    def __init__(self, enable_publish=True):
        self.sns = boto3.client('sns')
        self.enable_publish = enable_publish
        self.all_tickers = all_tickers

    def publish(self, arn, message):
        """Serialize a message and publish it to AWS SNS"""
        response = self.sns.publish(
            TargetArn=arn,
            Message=SNSMessage.serialize(message),
            MessageStructure=SNS.JSON,
        )

    def parse_publish_tickers(self, submission: SubmissionMessage):
        """Parse tickers from selftext and title fields to TickerMention messages and publish them to Kafka"""
        selftext_tickers = parse_tickers(
            submission.selftext, all_tickers=self.all_tickers
        )
        title_tickers = parse_tickers(submission.title, all_tickers=self.all_tickers)
        for ticker in selftext_tickers:
            mention = TickerMentionMessage(
                ticker,
                DS.REDDIT,
                PS.SUBMISSION_SELFTEXT,
                submission.submission_id,
                submission.created_utc,
            )
            if self.enable_publish:
                self.publish(SNS.TICKER_MENTIONS, mention)

        for ticker in title_tickers:
            mention = TickerMentionMessage(
                ticker,
                DS.REDDIT,
                PS.SUBMISSION_TITLE,
                submission.submission_id,
                submission.created_utc,
            )
            if self.enable_publish:
                self.publish(SNS.TICKER_MENTIONS, mention)

        has_tickers = len(selftext_tickers) > 0 or len(title_tickers) > 0
        return has_tickers

    def parse_publish_text(self, submission: SubmissionMessage):
        """Parse selftext and title fields to ScrapedPost messages and publish them to Kafka"""
        selftext_post = ScrapedPostMessage(
            submission.selftext,
            DS.REDDIT,
            PS.SUBMISSION_SELFTEXT,
            submission.submission_id,
            submission.created_utc,
        )
        if self.enable_publish:
            self.publish(SNS.SCRAPED_POSTS, selftext_post)

        title_post = ScrapedPostMessage(
            submission.title,
            DS.REDDIT,
            PS.SUBMISSION_TITLE,
            submission.submission_id,
            submission.created_utc,
        )
        if self.enable_publish:
            self.publish(SNS.SCRAPED_POSTS, title_post)

    def run(self, submission: SubmissionMessage):
        has_tickers = self.parse_publish_tickers(submission)
        if not has_tickers:
            return

        # Publish text fields as ScrapedPosts
        self.parse_publish_text(submission)

        return has_tickers