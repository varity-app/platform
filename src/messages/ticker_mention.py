from util.constants import TickerMentionsConstants as TMC

from . import Message

class TickerMentionMessage(Message):
    def __init__(self, stock_name: str, data_source: str, parent_source: str, parent_id: str, created_utc: str):
        self.stock_name = stock_name
        self.data_source = data_source
        self.parent_source = parent_source
        self.parent_id = parent_id
        self.created_utc = created_utc

    def serialize(self, as_str=True) -> str:
        response = dict(
            stock_name=self.stock_name,
            data_source=self.data_source,
            parent_source=self.parent_source,
            parent_id=self.parent_id,
            created_utc=self.created_utc,
        )

        if as_str:
            response = self.to_str(response)

        return response

    @classmethod
    def from_obj(cls, data: dict):
        for field in [TMC.STOCK_NAME, TMC.DATA_SOURCE, TMC.PARENT_SOURCE, TMC.PARENT_ID, TMC.CREATED_UTC]:
            assert field in data.keys()

        stock_name = data.get(TMC.STOCK_NAME)
        data_source = data.get(TMC.DATA_SOURCE)
        parent_source = data.get(TMC.PARENT_SOURCE)
        parent_id = data.get(TMC.PARENT_ID)
        created_utc = data.get(TMC.CREATED_UTC)
        
        return cls(
            stock_name,
            data_source,
            parent_source,
            parent_id,
            created_utc,
        )