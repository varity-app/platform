from confluent_kafka import Producer, Consumer
import logging
import json

from messages.reddit.submission import SubmissionMessage
from messages.scraped_post import ScrapedPostMessage
from messages.ticker_mention import TickerMentionMessage

from util.tickers import parse_tickers, all_tickers
from util.constants.scraping import DataSources as DS, ParentSources as PS
from util.constants.kafka import Config, Topics, Groups
from util.constants import MentionTypes

logger = logging.getLogger(__name__)


class SubmissionConsumer(object):
    def __init__(self, enable_publish=True):
        self.enable_publish = enable_publish

        if enable_publish:
            self.publisher = Producer(Config.OBJ)

        conf = {
            **Config.OBJ,
            'group.id': Groups.SUBMISSION_CONSUMERS,
            'auto.offset.reset': 'earliest'
        }
        self.consumer = Consumer(conf)

    def callback(self, err, msg):
        """Kafka publisher callback"""
        if err is not None:
            logger.error(f"Failed to deliver message: {err}")
        else:
            logger.debug(
                f"Produced record to topic {msg.topic()} partition [{msg.partition()}] @ offset {msg.offset()}")

    def publish(self, topic, message):
        """Serialize a message and publish it to Kafka"""
        self.publisher.produce(
            topic, value=message.serialize(), callback=self.callback)

    def parse_publish_tickers(self, submission: SubmissionMessage):
        """Parse tickers from selftext and title fields to TickerMention messages and publish them to Kafka"""
        selftext_tickers = parse_tickers(
            submission.selftext, all_tickers=all_tickers
        )
        title_tickers = parse_tickers(
            submission.title, all_tickers=all_tickers)

        for ticker in selftext_tickers:
            mention = TickerMentionMessage(
                ticker,
                DS.REDDIT,
                PS.SUBMISSION_SELFTEXT,
                submission.submission_id,
                submission.created_utc,
                MentionTypes.TICKER,
            )
            if self.enable_publish:
                self.publish(Topics.TICKER_MENTIONS, mention)

        for ticker in title_tickers:
            mention = TickerMentionMessage(
                ticker,
                DS.REDDIT,
                PS.SUBMISSION_TITLE,
                submission.submission_id,
                submission.created_utc,
                MentionTypes.TICKER,
            )
            if self.enable_publish:
                self.publish(Topics.TICKER_MENTIONS, mention)

        has_tickers = len(selftext_tickers) > 0 or len(title_tickers) > 0
        return has_tickers

    def parse_publish_text(self, submission: SubmissionMessage):
        """Parse selftext and title fields to ScrapedPost messages and publish them to Kafka"""
        if submission.selftext:
            selftext_post = ScrapedPostMessage(
                submission.selftext,
                DS.REDDIT,
                PS.SUBMISSION_SELFTEXT,
                submission.submission_id,
                submission.created_utc,
            )
            if self.enable_publish:
                self.publish(Topics.SCRAPED_POSTS, selftext_post)

        title_post = ScrapedPostMessage(
            submission.title,
            DS.REDDIT,
            PS.SUBMISSION_TITLE,
            submission.submission_id,
            submission.created_utc,
        )
        if self.enable_publish:
            self.publish(Topics.SCRAPED_POSTS, title_post)

    def run(self):
        self.consumer.subscribe([Topics.REDDIT_SUBMISSIONS])

        try:
            while True:
                msg = self.consumer.poll(1.0)
                if msg is None:
                    logger.debug("Waiting for message...")
                    continue
                elif msg.error():
                    logger.error(f"Error: {msg.error()}")
                    continue

                # Load payload
                data = json.loads(msg.value())
                submission = SubmissionMessage.from_obj(data)

                # Check for tickers
                has_tickers = self.parse_publish_tickers(submission)
                if not has_tickers:
                    continue

                # Publish text fields as ScrapedPosts
                self.parse_publish_text(submission)

                self.publisher.flush()

        except KeyboardInterrupt:
            logger.error("Keyboard interrupt")
            pass
        finally:
            self.consumer.close()
