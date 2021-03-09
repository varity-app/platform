import logging

from consumers.sentiment import SentimentEstimator

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def main():
    consumer = SentimentEstimator()

    logger.info("Starting to watch Kafka...")
    consumer.run()


if __name__ == "__main__":
    main()