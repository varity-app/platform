import json
import base64

from . import Message


class SNSMessage(object):

    @staticmethod
    def serialize(message: Message):
        encoded_message = message.serialize(as_str=True)
        encoded_message = base64.b64encode(encoded_message.encode("utf-8"))
        encoded_message = encoded_message.decode("utf-8")

        response = dict(
            default=encoded_message
        )

        response = json.dumps(response)

        return response

    @staticmethod
    def parse_message(record):
        response = record['body']
        response = json.loads(response)
        response = response['Message']
        response = base64.decodebytes(response.encode("utf-8"))
        response = json.loads(response)

        return response
