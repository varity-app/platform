"""
Declare Kafka message class for logs
"""

from util.constants.logging import Constants

from . import Message


class LogMessage(Message):
    """Kafka message repesenting a log entry"""

    def __init__(
        self, timestamp: str, source_type: str, source: str, message: str
    ) -> None:
        self.timestamp = timestamp
        self.source_type = source_type
        self.source = source
        self.message = message

    def to_obj(self) -> dict:
        """Serialize the class as a dict"""
        obj = dict(
            timestamp=self.timestamp,
            source_type=self.source_type,
            source=self.source,
            message=self.message,
        )

        return obj

    @classmethod
    def from_obj(cls: LogMessage, data: dict) -> LogMessage:
        """Create instance of LogMessage from a dictionary"""
        fields = [
            Constants.TIMESTAMP,
            Constants.SOURCE_TYPE,
            Constants.SOURCE,
            Constants.MESSAGE,
        ]

        cls.assert_has_fields(data, fields)

        timestamp = data[Constants.TIMESTAMP]
        source_type = data[Constants.SOURCE_TYPE]
        source = data[Constants.SOURCE]
        message = data[Constants.MESSAGE]

        return cls(timestamp, source_type, source, message)
