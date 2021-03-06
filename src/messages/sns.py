import json

from . import Message


class SNSMessage(object):

    @staticmethod
    def serialize(message: Message):
        response = dict(
            default=message.serialize(as_str=True)
        )

        response = json.dumps(response)

        return response
