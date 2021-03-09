import logging

from consumers.reddit.comment import CommentConsumer

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def main():
    consumer = CommentConsumer()

    logger.info("Starting to watch Kafka...")
    consumer.run()


if __name__ == "__main__":
    main()