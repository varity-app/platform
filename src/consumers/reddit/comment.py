from confluent_kafka import Producer, Consumer
import logging
import json

from messages.reddit.comment import CommentMessage
from messages.scraped_post import ScrapedPostMessage
from messages.ticker_mention import TickerMentionMessage

from util.tickers import parse_tickers, all_tickers
from util.constants.scraping import DataSources as DS, ParentSources as PS
from util.constants.kafka import Config, Topics, Groups
from util.constants import MentionTypes

logger = logging.getLogger(__name__)


class CommentConsumer(object):
    def __init__(self, enable_publish=True):
        self.enable_publish = enable_publish

        if enable_publish:
            self.publisher = Producer(Config.OBJ)

        conf = {
            **Config.OBJ,
            'group.id': Groups.COMMENT_CONSUMERS,
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
        self.publisher.produce(topic, value=message.serialize(), callback=self.callback)

    def parse_publish_tickers(self, comment: CommentMessage):
        """Parse tickers from body to TickerMention messages and publish to Kafka"""
        body_tickers = parse_tickers(
            comment.body, all_tickers=all_tickers
        )
        for ticker in body_tickers:
            mention = TickerMentionMessage(
                ticker,
                DS.REDDIT,
                PS.COMMENT_BODY,
                comment.comment_id,
                comment.created_utc,
                MentionTypes.TICKER,
            )
            if self.enable_publish:
                self.publish(Topics.TICKER_MENTIONS, mention)

        has_tickers = len(body_tickers) > 0
        return has_tickers

    def parse_publish_text(self, comment: CommentMessage):
        """Parse body field to ScrapedPost message and publish to Kafka"""
        body_post = ScrapedPostMessage(
            comment.body,
            DS.REDDIT,
            PS.COMMENT_BODY,
            comment.comment_id,
            comment.created_utc,
        )
        if self.enable_publish:
            self.publish(Topics.SCRAPED_POSTS, body_post)

    def run(self):
        self.consumer.subscribe([Topics.REDDIT_COMMENTS])

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
                comment = CommentMessage.from_obj(data)

                # Check for tickers
                has_tickers = self.parse_publish_tickers(comment)
                if not has_tickers:
                    continue

                # Publish text fields as ScrapedPosts
                self.parse_publish_text(comment)

                self.publisher.flush()

        except KeyboardInterrupt:
            logger.error("Keyboard interrupt")
            pass
        finally:
            self.consumer.close()
