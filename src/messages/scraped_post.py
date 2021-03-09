from util.constants import ScrapedPostConstants as SPC
from . import Message


class ScrapedPostMessage(Message):
    def __init__(self, text: str, data_source: str, parent_source: str, parent_id: str, timestamp: str):
        self.text = text
        self.data_source = data_source
        self.parent_source = parent_source
        self.parent_id = parent_id
        self.timestamp = timestamp

    
    def serialize(self, as_str=True) -> dict:
        response = dict(
            text=self.text,
            data_source=self.data_source,
            parent_source=self.parent_source,
            parent_id=self.parent_id,
            timestamp=self.timestamp,
        )

        if as_str:
            response = self.to_str(response)

        return response

    @classmethod
    def from_obj(cls, data: dict):
        for field in [SPC.TEXT, SPC.DATA_SOURCE, SPC.PARENT_SOURCE, SPC.PARENT_ID, SPC.TIMESTAMP]:
            assert field in data.keys()

        text = data.get(SPC.TEXT)
        data_source = data.get(SPC.DATA_SOURCE)
        parent_source = data.get(SPC.PARENT_SOURCE)
        parent_id = data.get(SPC.PARENT_ID)
        timestamp = data.get(SPC.TIMESTAMP)

        return cls(text, data_source, parent_source, parent_id, timestamp)

        