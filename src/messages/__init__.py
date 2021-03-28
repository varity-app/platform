"""
Declare the base helper class for Kafka Messages.
All other Kafka Message helper classes should inherit from this one.
"""

from abc import ABCMeta, abstractmethod
from typing import Union, List
import json


class Message:
    """Base helper class for a Kafka Message"""

    __metaclass__ = ABCMeta

    @abstractmethod
    def to_obj(self):
        """Serialize a Message to a dict"""
        return dict()

    def serialize(self, as_str=True) -> Union[dict, str]:
        """Serialize a Message to either a dict or a str"""
        response = self.to_obj()

        if as_str:
            response = json.dumps(response)

        return response

    @staticmethod
    def get(value, default=""):
        """
        Assign a value to a default if it is None.
        This is done to ensure that a scraped null value has a reasonable default.
        """
        if value is None:
            return default

        return value

    @staticmethod
    def assert_has_fields(obj: dict, fields: List[str]) -> None:
        """Assert that an object has all required fields"""
        for field in fields:
            assert field in obj.keys()
