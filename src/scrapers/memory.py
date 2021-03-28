"""
Declare the Memory helper class
"""


class Memory:
    """Helper class that helps a reddit scraper determine if a post has been already scraped"""

    def __init__(self, memory_size: int) -> None:
        self.memory_size = memory_size
        self.store = []
        self.cursor = 0

    def add(self, data: str):
        """Add an item to memory"""
        if len(self.store) >= self.memory_size:
            self.store.pop()
            self.store.insert(self.cursor, data)
            self.cursor += 1
        else:
            self.store.append(data)

    def reset_cursor(self):
        """Reset the memory cursor to the beginning"""
        self.cursor = 0

    def contains(self, data: str):
        """Return whether an item is in memory"""
        return data in self.store
