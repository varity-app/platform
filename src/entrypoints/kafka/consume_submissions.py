import logging

from consumers.reddit.submission import SubmissionConsumer

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def main():
    consumer = SubmissionConsumer()

    logger.info("Starting to watch Kafka...")
    consumer.run()


if __name__ == "__main__":
    main()