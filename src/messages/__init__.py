from abc import ABCMeta, abstractmethod
import json


class Message(object):
    __metaclass__ = ABCMeta

    @abstractmethod
    def serialize(self, as_str=True):
        response = dict()

        if as_str:
            response = self.to_str(response)

        return response

    @staticmethod
    def to_str(obj):
        return json.dumps(obj)
