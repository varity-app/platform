from confluent_kafka import Producer
import json
import logging

from util.constants.kafka import Config, Topics

logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)


def callback(err, msg):
    if err is not None:
        logger.error(f"Failed to deliver message: {err}")
    else:
        logger.info(
            f"Produced record to topic {msg.topic()} partition [{msg.partition()}] @ offset {msg.offset()}")


def main():
    logger.info("Testing Confluence Kafka")
    producer = Producer(Config.OBJ)

    producer.produce(Topics.REDDIT_SUBMISSIONS,
                     value=json.dumps(dict(message="yeet")),
                     on_delivery=callback)
    producer.flush()


if __name__ == '__main__':
    main()
