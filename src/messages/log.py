"""
Declare Kafka message class for logs
"""

from typing import Union, Dict

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

    def serialize(self, as_str=True) -> Union[str, Dict[str, str]]:
        """Serialize message to dict or str"""

        response = dict(
            timestamp=self.timestamp,
            source_type=self.source_type,
            source=self.source,
            message=self.message,
        )

        if as_str:
            response = self.to_str(response)

        return response

    @classmethod
    def from_obj(cls, data: dict):
        """Create instance of LogMessage from a dictionary"""

        for field in [
            Constants.TIMESTAMP,
            Constants.SOURCE_TYPE,
            Constants.SOURCE,
            Constants.MESSAGE,
        ]:
            assert field in data.keys()

        timestamp = data[Constants.TIMESTAMP]
        source_type = data[Constants.SOURCE_TYPE]
        source = data[Constants.SOURCE]
        message = data[Constants.MESSAGE]

        return cls(timestamp, source_type, source, message)
