import faust

class ScrapedPost(faust.Record):
    text: str
    data_source: str
    parent_source: str
    parent_id: str
    timestamp: str