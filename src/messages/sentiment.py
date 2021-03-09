from . import Message

from util.constants import ScrapedPostConstants as SPC
from util.constants.sentiment import Constants as C


class SentimentMessage(Message):
    def __init__(self, data_source: str, parent_source: str, parent_id: str, estimator: str, compound: float, positive: float, neutral: float, negative: float):
        self.data_source = data_source
        self.parent_source = parent_source
        self.parent_id = parent_id
        self.compound = compound
        self.positive = positive
        self.neutral = neutral
        self.negative = negative
        self.estimator = estimator

    def serialize(self, as_str=True) -> dict:
        response = dict(
            data_source=self.data_source,
            parent_source=self.parent_source,
            parent_id=self.parent_id,
            compound=self.compound,
            positive=self.positive,
            negative=self.negative,
            neutral=self.neutral,
            estimator=self.estimator,
        )

        if as_str:
            response = self.to_str(response)

        return response

    @classmethod
    def from_obj(cls, data: dict):
        for field in [SPC.DATA_SOURCE, SPC.PARENT_SOURCE, SPC.PARENT_ID, C.ESTIMATOR, C.COMPOUND, C.POSITIVE, C.NEUTRAL, C.NEGATIVE]:
            assert field in data.keys()

        data_source = data.get(SPC.DATA_SOURCE)
        parent_source = data.get(SPC.PARENT_SOURCE)
        parent_id = data.get(SPC.PARENT_ID)
        estimator = data.get(C.ESTIMATOR)
        compound = data.get(C.COMPOUND)
        positive = data.get(C.POSITIVE)
        neutral = data.get(C.NEUTRAL)
        negative = data.get(C.NEGATIVE)

        return cls(data_source, parent_source, parent_id, estimator, compound, positive, neutral, negative)
