"""
Faust data models for ScrapedPosts
"""

import faust


class ScrapedPost(faust.Record):
    """Faust model for ScrapedPosts"""

    text: str
    data_source: str
    parent_source: str
    parent_id: str
    timestamp: str
