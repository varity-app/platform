"""
Custom handler for pushing logs to Kafka
"""

import logging
from datetime import datetime
from confluent_kafka import Producer

from messages.log import LogMessage
from util.constants.kafka import Config, Topics
from util.constants.logging import SourceTypes


class KafkaHandler(logging.StreamHandler):
    """Handler for pushing logs to Kafka"""

    def __init__(self, source: str) -> None:
        logging.StreamHandler.__init__(self)
        self.producer = Producer(Config.OBJ)

        self.source = source

    def emit(self, record):
        """Emit a log to kafka"""

        # Format Log
        msg = self.format(record)
        log = LogMessage(
            datetime.now().isoformat(),
            SourceTypes.PYTHON,
            self.source,
            msg,
        )

        # Publish to Kafka
        self.producer.produce(Topics.LOGS, value=log.serialize())
        self.producer.flush()
