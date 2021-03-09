from confluent_kafka import Producer, Consumer
from nltk.sentiment.vader import SentimentIntensityAnalyzer
import json
import logging

from messages.scraped_post import ScrapedPostMessage
from messages.sentiment import SentimentMessage
from util.constants.kafka import Config, Topics, Groups
from util.constants.sentiment import Estimators

logger = logging.getLogger(__name__)


class SentimentEstimator(object):
    def __init__(self, enable_publish=True):
        self.enable_publish = enable_publish

        if enable_publish:
            self.publisher = Producer(Config.OBJ)

        conf = {
            **Config.OBJ,
            'group.id': Groups.SENTIMENT_ESTIMATORS,
            'auto.offset.reset': 'earliest'
        }
        self.consumer = Consumer(conf)

        self.model = SentimentIntensityAnalyzer()

    def callback(self, err, msg):
        """Kafka publisher callback"""
        if err is not None:
            logger.error(f"Failed to deliver message: {err}")
        else:
            logger.debug(
                f"Produced record to topic {msg.topic()} partition [{msg.partition()}] @ offset {msg.offset()}")

    def run(self):
        self.consumer.subscribe([Topics.SCRAPED_POSTS])

        try:
            while True:
                msg = self.consumer.poll(1.0)
                if msg is None:
                    logger.debug("Waiting for message...")
                    continue
                elif msg.error():
                    logger.error(f"Error: {msg.error()}")
                    continue

                # Parse ScrapedPostMessage
                data = json.loads(msg.value())
                post = ScrapedPostMessage.from_obj(data)

                # Estimate sentiment
                scores = self.model.polarity_scores(post.text)

                # Parse object
                sentiment = SentimentMessage(
                    post.data_source,
                    post.parent_source,
                    post.parent_id,
                    Estimators.VADER,
                    scores['compound'],
                    scores['pos'],
                    scores['neu'],
                    scores['neg'],
                )

                # Publish to Kafka
                if self.enable_publish:
                    self.publisher.produce(Topics.POST_SENTIMENT,
                                           value=sentiment.serialize(),
                                           callback=self.callback)

                self.publisher.flush()

        except KeyboardInterrupt:
            logger.error("Keyboard interrupt")
            pass
        finally:
            self.consumer.close()
