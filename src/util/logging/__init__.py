"""
Helper module for a custom logger
"""

import logging

from .handler import KafkaHandler

logging.basicConfig(
    format="%(name)s:%(levelname)s:%(message)s",
    level=logging.INFO,
)


def get_logger(source, enable_kafka=True):
    """Get a custom logger optionally wrapped in a KafkaHandler"""

    logger = logging.getLogger(source)

    if enable_kafka:
        handler = KafkaHandler(source)

        # Set formatter
        formatter = logging.Formatter("%(levelname)s:%(message)s")
        handler.setFormatter(formatter)

        # Add handler to logger
        logger.setLevel(logging.INFO)
        logger.addHandler(handler)

    return logger
